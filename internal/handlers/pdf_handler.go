package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// PDFHandler handles PDF generation requests with quota management
type PDFHandler struct{}

// NewPDFHandler creates a new PDF handler
func NewPDFHandler() *PDFHandler {
	return &PDFHandler{}
}

// PDFRequest represents a PDF generation request
type PDFRequest struct {
	Type    string                 `json:"type" validate:"required,oneof=nutrition workout health recipes complete"`
	Format  string                 `json:"format,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// PDFResponse represents a PDF generation response
type PDFResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Usage       int    `json:"usage"`
	Limit       int    `json:"limit"`
	DownloadURL string `json:"download_url,omitempty"`
}

// getUserPlan extracts user plan from cookies
func (h *PDFHandler) getUserPlan(c echo.Context) (plan string, shared bool) {
	planCookie, err := c.Cookie("plan")
	if err == nil {
		plan = planCookie.Value
	} else {
		plan = "free"
	}

	sharedCookie, err := c.Cookie("shared")
	shared = err == nil && sharedCookie.Value == "yes"

	return plan, shared
}

// getPDFQuota determines the PDF quota based on plan and shared status
func (h *PDFHandler) getPDFQuota(plan string, shared bool) int {
	switch plan {
	case "pro":
		return 50
	case "premium", "lifetime":
		return 999999 // Unlimited
	case "free":
		if shared {
			return 11
		}
		return 3
	default:
		return 3
	}
}

// getCurrentUsage gets current PDF usage from cookies
func (h *PDFHandler) getCurrentUsage(c echo.Context) int {
	usageCookie, err := c.Cookie("pdfUsage")
	if err != nil {
		return 0
	}

	usage, err := strconv.Atoi(usageCookie.Value)
	if err != nil {
		return 0
	}

	return usage
}

// updateUsage increments the PDF usage counter
func (h *PDFHandler) updateUsage(c echo.Context, currentUsage int) {
	newUsage := currentUsage + 1
	cookie := &http.Cookie{
		Name:     "pdfUsage",
		Value:    strconv.Itoa(newUsage),
		Path:     "/",
		Expires:  time.Now().AddDate(0, 1, 0), // 1 month
		HttpOnly: false,                       // Allow JavaScript access for quota display
		Secure:   false,                       // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
}

// GeneratePDF handles PDF generation requests
func (h *PDFHandler) GeneratePDF(c echo.Context) error {
	var req PDFRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, PDFResponse{
			Success: false,
			Message: "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, PDFResponse{
			Success: false,
			Message: "Invalid PDF type. Must be one of: nutrition, workout, health, recipes, complete",
		})
	}

	// Get user plan and quota
	plan, shared := h.getUserPlan(c)
	quota := h.getPDFQuota(plan, shared)
	currentUsage := h.getCurrentUsage(c)

	// Check quota limits (except for unlimited plans)
	if quota != 999999 && currentUsage >= quota {
		return c.JSON(http.StatusTooManyRequests, PDFResponse{
			Success: false,
			Message: fmt.Sprintf("PDF quota exceeded. You've used %d/%d PDFs this month. Upgrade your plan for more downloads.", currentUsage, quota),
			Usage:   currentUsage,
			Limit:   quota,
		})
	}

	// Generate PDF content based on type
	content, err := h.generatePDFContent(req.Type, req.Options)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, PDFResponse{
			Success: false,
			Message: "Failed to generate PDF content",
		})
	}

	// Update usage counter
	h.updateUsage(c, currentUsage)

	// Return PDF content as text (in production, this would be actual PDF bytes)
	fileName := fmt.Sprintf("doctorhealthy1-%s-report-%s.txt", req.Type, time.Now().Format("2006-01-02"))

	c.Response().Header().Set("Content-Type", "text/plain")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return c.String(http.StatusOK, content)
}

// GetQuotaStatus returns current PDF usage and limits
func (h *PDFHandler) GetQuotaStatus(c echo.Context) error {
	plan, shared := h.getUserPlan(c)
	quota := h.getPDFQuota(plan, shared)
	currentUsage := h.getCurrentUsage(c)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"usage":  currentUsage,
		"limit":  quota,
		"plan":   plan,
		"shared": shared,
		"remaining": func() int {
			if quota == 999999 {
				return 999999
			}
			remaining := quota - currentUsage
			if remaining < 0 {
				return 0
			}
			return remaining
		}(),
	})
}

// generatePDFContent creates content based on PDF type
func (h *PDFHandler) generatePDFContent(pdfType string, options map[string]interface{}) (string, error) {
	timestamp := time.Now().Format("January 2, 2006 at 3:04 PM")

	templates := map[string]string{
		"nutrition": fmt.Sprintf(`DoctorHealthy1 - Nutrition Report
Generated: %s

Your Personalized Nutrition Plan
===============================

Based on your profile and health goals, here's your customized nutrition guide:

Daily Caloric Needs: 2,200 calories

Macronutrient Breakdown:
- Proteins: 25%% (550 cal, 137g)
- Carbohydrates: 45%% (990 cal, 247g) 
- Fats: 30%% (660 cal, 73g)

Recommended Foods:
✓ Lean proteins: Chicken breast, fish, legumes, Greek yogurt
✓ Complex carbs: Quinoa, brown rice, sweet potatoes, oats
✓ Healthy fats: Avocados, nuts, olive oil, salmon

Meal Timing Recommendations:
- Breakfast: 7-8 AM (400 calories)
- Mid-morning snack: 10 AM (150 calories)
- Lunch: 12-1 PM (600 calories)
- Afternoon snack: 3 PM (200 calories)
- Dinner: 6-7 PM (550 calories)
- Evening snack: 8 PM (150 calories)

Weekly Meal Planning Tips:
1. Prep proteins in bulk on Sundays
2. Keep healthy snacks readily available
3. Stay hydrated - aim for 8-10 glasses of water daily
4. Track your intake for the first 2 weeks

Next Review: %s

Generated by DoctorHealthy1 AI System
Visit doctorhealthy1.com for more personalized recommendations`, timestamp, time.Now().AddDate(0, 1, 0).Format("January 2, 2006")),

		"workout": fmt.Sprintf(`DoctorHealthy1 - Workout Plan
Generated: %s

Your Personalized 4-Week Fitness Program
=======================================

Current Fitness Level: Intermediate
Goal: Strength Building & Fat Loss

Week 1-2: Foundation Phase
------------------------

MONDAY - Upper Body Strength:
• Push-ups: 3 sets x 12 reps
• Pull-ups (assisted if needed): 3 sets x 8 reps
• Dumbbell rows: 3 sets x 10 reps
• Plank: 3 sets x 45 seconds
• Rest: 60-90 seconds between sets

WEDNESDAY - Lower Body Power:
• Squats: 4 sets x 15 reps
• Lunges: 3 sets x 10 each leg
• Glute bridges: 3 sets x 15 reps
• Calf raises: 3 sets x 20 reps
• Wall sit: 3 sets x 30 seconds

FRIDAY - Full Body Circuit:
• Burpees: 3 sets x 8 reps
• Mountain climbers: 3 sets x 20 reps
• Jumping jacks: 3 sets x 30 seconds
• High knees: 3 sets x 20 reps
• Cool down: 10 minutes stretching

Week 3-4: Strength Phase
-----------------------

MONDAY - Advanced Upper Body:
• Push-ups: 4 sets x 15 reps
• Pull-ups: 4 sets x 10 reps
• Pike push-ups: 3 sets x 8 reps
• Plank to downward dog: 3 sets x 10 reps
• Side planks: 3 sets x 30 seconds each side

WEDNESDAY - Advanced Lower Body:
• Jump squats: 4 sets x 12 reps
• Bulgarian split squats: 3 sets x 8 each leg
• Single-leg glute bridges: 3 sets x 10 each leg
• Pistol squat progression: 3 sets x 5 each leg

FRIDAY - HIIT Circuit:
• 20 seconds work, 10 seconds rest
• 8 rounds of mixed exercises
• Focus on maximum intensity

Recovery & Nutrition:
- Sleep: 7-9 hours nightly
- Protein: 1.6-2.2g per kg body weight
- Hydration: 3-4 liters daily
- Rest days: Light walking or yoga

Progress Tracking:
□ Week 1 measurements
□ Week 2 progress photos
□ Week 3 strength assessment
□ Week 4 final evaluation

Next Program Update: %s

Generated by DoctorHealthy1 AI System`, timestamp, time.Now().AddDate(0, 1, 0).Format("January 2, 2006")),

		"health": fmt.Sprintf(`DoctorHealthy1 - Health Assessment Report
Generated: %s

Comprehensive Health Analysis
===========================

Overall Health Score: 85/100 (Excellent)

Key Health Metrics:
• BMI: 23.5 (Normal range: 18.5-24.9)
• Body Fat Percentage: 15%% (Excellent)
• Muscle Mass Index: 42%% (Good)
• Hydration Level: 68%% (Adequate)
• Sleep Quality: 7.2/10 (Good)
• Stress Level: 4/10 (Low-Moderate)

Cardiovascular Health:
✓ Resting heart rate: 65 bpm (Good)
✓ Blood pressure: 120/80 (Normal)
✓ Cardio fitness: Above average
✓ Recovery rate: Excellent

Nutritional Status:
✓ Vitamin D: Adequate
✓ Iron levels: Normal
✓ B12: Good
⚠ Omega-3: Could be improved
⚠ Fiber intake: Below recommended

Physical Assessment:
✓ Flexibility: Good
✓ Core strength: Above average
✓ Balance: Excellent
⚠ Upper body strength: Needs improvement

Health Recommendations:
1. IMMEDIATE (1-2 weeks):
   - Increase omega-3 rich foods (salmon, walnuts)
   - Add 5g more fiber daily
   - Practice deep breathing exercises daily

2. SHORT-TERM (1-3 months):
   - Incorporate 2x weekly strength training
   - Improve sleep hygiene (aim for 8 hours)
   - Regular health check-ups

3. LONG-TERM (3-12 months):
   - Maintain current cardio routine
   - Consider comprehensive blood panel
   - Develop stress management techniques

Risk Assessment:
• Low Risk: Heart disease, diabetes, hypertension
• Medium Risk: Stress-related disorders
• Action Needed: Posture improvement, flexibility work

Lifestyle Factors:
✓ Non-smoker
✓ Moderate alcohol consumption
✓ Regular exercise routine
⚠ High screen time (consider breaks)
⚠ Irregular meal timing

Health Goals Progress:
• Weight management: On track ✓
• Fitness improvement: Ahead of schedule ✓
• Energy levels: Significantly improved ✓
• Sleep quality: Needs attention ⚠

Next Health Review: %s
Recommended Follow-up: 6 weeks

Generated by DoctorHealthy1 AI System
For personalized health coaching, visit doctorhealthy1.com`, timestamp, time.Now().AddDate(0, 1, 15).Format("January 2, 2006")),

		"recipes": fmt.Sprintf(`DoctorHealthy1 - Recipe Collection
Generated: %s

Your Personalized Healthy Recipe Book
===================================

BREAKFAST RECIPES
================

1. Protein Power Bowl
Prep time: 10 minutes | Serves: 1 | Calories: 420

Ingredients:
• 2 large eggs
• 1/2 avocado, sliced
• 1 cup fresh spinach
• 1/4 cup cooked quinoa
• 1 tbsp olive oil
• Salt, pepper, paprika to taste

Instructions:
1. Scramble eggs in olive oil
2. Arrange spinach, quinoa in bowl
3. Top with eggs and avocado
4. Season and serve immediately

Nutrition: Protein 22g | Carbs 18g | Fat 28g

2. Green Power Smoothie
Prep time: 5 minutes | Serves: 1 | Calories: 280

Ingredients:
• 1 banana
• 1 cup spinach
• 1/2 cup unsweetened almond milk
• 1 tbsp almond butter
• 1 tsp chia seeds
• Ice cubes

Instructions:
1. Blend all ingredients until smooth
2. Add ice for desired thickness
3. Serve immediately

Nutrition: Protein 8g | Carbs 32g | Fat 14g

LUNCH RECIPES
============

3. Mediterranean Power Bowl
Prep time: 15 minutes | Serves: 1 | Calories: 520

Ingredients:
• 4 oz grilled chicken breast
• 1/2 cup cooked quinoa
• 1/4 cup cherry tomatoes
• 1/4 cucumber, diced
• 2 tbsp hummus
• 1 tbsp olive oil
• Fresh herbs (parsley, mint)

Instructions:
1. Season and grill chicken
2. Arrange quinoa in bowl
3. Top with vegetables and chicken
4. Drizzle with olive oil
5. Serve with hummus

Nutrition: Protein 35g | Carbs 45g | Fat 18g

4. Asian-Inspired Lettuce Wraps
Prep time: 20 minutes | Serves: 2 | Calories: 320 per serving

Ingredients:
• 6 oz ground turkey (lean)
• 1 head butter lettuce
• 1 cup shredded cabbage
• 1 carrot, julienned
• 2 tbsp low-sodium soy sauce
• 1 tbsp sesame oil
• 1 tsp fresh ginger, minced

Instructions:
1. Cook turkey with ginger and soy sauce
2. Prepare vegetable filling
3. Assemble in lettuce cups
4. Drizzle with sesame oil

Nutrition: Protein 28g | Carbs 12g | Fat 16g

DINNER RECIPES
=============

5. Herb-Crusted Salmon
Prep time: 25 minutes | Serves: 1 | Calories: 580

Ingredients:
• 5 oz salmon fillet
• 1 medium sweet potato
• 1 cup broccoli florets
• 2 tbsp olive oil
• Fresh dill, lemon
• Garlic powder, paprika

Instructions:
1. Roast sweet potato at 400°F for 25 min
2. Season salmon with herbs
3. Pan-sear salmon 4 min each side
4. Steam broccoli until tender
5. Serve with lemon wedges

Nutrition: Protein 42g | Carbs 35g | Fat 28g

HEALTHY SNACKS
=============

6. Energy Balls
Prep time: 10 minutes | Makes: 12 balls | Calories: 95 each

Ingredients:
• 1 cup dates, pitted
• 1/2 cup almonds
• 2 tbsp cocoa powder
• 1 tbsp coconut oil
• Pinch of sea salt

Instructions:
1. Process dates and almonds
2. Add cocoa and coconut oil
3. Roll into balls
4. Refrigerate 30 minutes

Nutrition per ball: Protein 2g | Carbs 12g | Fat 5g

MEAL PREP TIPS:
• Batch cook quinoa and proteins on Sundays
• Pre-cut vegetables for quick assembly
• Store dressings separately
• Freeze smoothie ingredients in portioned bags

SHOPPING LIST FOR THIS WEEK:
Proteins: Eggs, chicken breast, salmon, ground turkey
Vegetables: Spinach, avocado, broccoli, sweet potatoes
Pantry: Quinoa, almond butter, olive oil, spices

Next Recipe Collection: %s

Generated by DoctorHealthy1 AI System
Discover more recipes at doctorhealthy1.com`, timestamp, time.Now().AddDate(0, 0, 7).Format("January 2, 2006")),

		"complete": fmt.Sprintf(`DoctorHealthy1 - Complete Health Report
Generated: %s

COMPREHENSIVE WELLNESS PROFILE
=============================

Personal Information:
• Report Date: %s
• Plan: Premium
• Status: Active Member
• Next Review: %s

EXECUTIVE SUMMARY
================

Overall Health Score: 85/100 (Excellent)

You're making excellent progress on your health journey! Your commitment to 
nutrition and fitness is showing measurable results. Focus areas for the 
next month include improving sleep quality and increasing upper body strength.

KEY ACHIEVEMENTS THIS MONTH:
✓ Maintained target weight range
✓ Improved cardiovascular endurance by 12%%
✓ Increased daily protein intake to optimal levels
✓ Consistent workout routine (5 days/week)

HEALTH METRICS ANALYSIS
======================

Biometric Data:
• Height: 5'9" (175 cm)
• Weight: 165 lbs (75 kg) - Target range: 160-170 lbs ✓
• BMI: 23.5 (Normal) ✓
• Body Fat: 15%% (Excellent) ✓
• Muscle Mass: 42%% (Good, improving) ↗
• Hydration: 68%% (Adequate) ⚠

Cardiovascular Health:
• Resting HR: 65 bpm (Good) ✓
• Max HR: 185 bpm (Excellent) ✓
• Blood Pressure: 120/80 (Optimal) ✓
• VO2 Max: 45 ml/kg/min (Above Average) ✓

NUTRITION ANALYSIS
=================

Daily Averages (Last 30 days):
• Calories: 2,180 (Target: 2,200) ✓
• Protein: 135g (25%%) ✓
• Carbohydrates: 245g (45%%) ✓
• Fats: 72g (30%%) ✓
• Fiber: 28g (Good) ✓
• Water: 2.1L (Increase to 2.5L) ⚠

Micronutrient Status:
✓ Vitamin D: Adequate
✓ Iron: Normal range
✓ B12: Excellent
✓ Calcium: Good
⚠ Omega-3: Could improve
⚠ Magnesium: Below optimal

Top Nutrition Wins:
1. Consistent protein timing around workouts
2. Increased vegetable intake by 30%%
3. Reduced processed food consumption
4. Better meal prep consistency

Areas for Improvement:
1. Add 2 servings of fatty fish weekly
2. Include more magnesium-rich foods
3. Increase water intake by 400ml daily

FITNESS PROGRESS REPORT
======================

Workout Consistency: 87%% (Excellent)
• Target: 5 sessions/week
• Achieved: 4.3 sessions/week average
• Missed sessions: Usually due to travel

Strength Gains (8-week comparison):
• Push-ups: 12 → 20 reps (+67%%) ✓
• Pull-ups: 3 → 7 reps (+133%%) ✓
• Plank hold: 45s → 90s (+100%%) ✓
• Squat: 15 → 25 reps (+67%%) ✓

Cardio Improvements:
• 5K time: 28:30 → 25:45 (-10%%) ✓
• Resting HR: 72 → 65 bpm ✓
• Recovery time: Significantly improved ✓

Current Program: Intermediate Strength + Cardio
Next Phase: Advanced strength with periodization

LIFESTYLE FACTORS
================

Sleep Analysis:
• Average: 7.2 hours/night (Target: 8 hours) ⚠
• Quality: 7.8/10 (Good) ✓
• Wake time consistency: Excellent ✓
• Recommendation: Earlier bedtime routine

Stress Management:
• Level: 4/10 (Low-moderate) ✓
• Coping strategies: Exercise, meditation ✓
• Work-life balance: Good ✓
• Recommendation: Continue current practices

Screen Time:
• Daily average: 6.8 hours ⚠
• Recommendation: Implement 20-20-20 rule
• Blue light: Use filters after 8 PM

SUPPLEMENT RECOMMENDATIONS
========================

Current Stack (Recommended):
✓ Vitamin D3: 2000 IU daily
✓ Omega-3: 1000mg EPA/DHA daily  
✓ Magnesium Glycinate: 200mg before bed
✓ Probiotic: 10 billion CFU daily

Optional Additions:
• Creatine: 3-5g daily (for strength gains)
• Vitamin B Complex: If energy levels drop
• Zinc: 15mg daily (immune support)

GOAL TRACKING
============

3-Month Goals Progress:
1. Lose 10 lbs: 8 lbs down ✓ (80%% complete)
2. Run 5K under 25 min: 25:45 ✓ (Achieved!)
3. 10 pull-ups: 7 pull-ups ↗ (70%% complete)
4. Improve sleep: 6.5 → 7.2 hrs ↗ (Improving)

Next 3-Month Goals:
1. Maintain current weight (160-170 lbs)
2. Increase pull-ups to 12 reps
3. Add 2 yoga sessions weekly
4. Achieve 8 hours average sleep

HEALTH RISK ASSESSMENT
=====================

Current Risk Levels:
• Cardiovascular Disease: Low ✓
• Type 2 Diabetes: Very Low ✓
• Metabolic Syndrome: Very Low ✓
• Injury Risk: Low-Moderate ⚠

Protective Factors:
✓ Regular exercise routine
✓ Healthy BMI and body composition
✓ Non-smoker
✓ Moderate alcohol consumption
✓ Strong social support system

Areas Requiring Attention:
⚠ Improve flexibility and mobility work
⚠ Address minor postural imbalances
⚠ Enhance recovery protocols

ACTION PLAN - NEXT 30 DAYS
==========================

Week 1-2: Foundation
• Continue current workout routine
• Add 15 minutes daily mobility work
• Increase water to 2.5L daily
• Track sleep patterns

Week 3-4: Enhancement
• Add omega-3 rich foods (salmon, walnuts)
• Implement blue light blocking routine
• Progressive overload in strength exercises
• Schedule comprehensive blood panel

Monthly Targets:
□ Complete 20 workouts
□ Achieve 8-hour sleep average
□ Increase pull-ups to 8 reps
□ Maintain nutrition consistency above 85%%

RESOURCES & SUPPORT
==================

Your DoctorHealthy1 Tools:
• Nutrition tracker with AI recommendations
• Personalized workout programs
• Progress photos and measurements
• 24/7 health coaching chat
• Recipe database with 1000+ options

Community Support:
• Monthly group challenges
• Expert Q&A sessions
• Member success stories
• Recipe sharing community

Next Steps:
1. Review this report with your health coach
2. Update goals based on current progress
3. Schedule follow-up assessment in 4 weeks
4. Continue daily habit tracking

CONTACT INFORMATION
==================

DoctorHealthy1 Support Team
Email: support@doctorhealthy1.com
Website: doctorhealthy1.com
Emergency Health Line: Available 24/7

Remember: Consistency beats perfection. You're making excellent progress!

Report End
Generated by DoctorHealthy1 AI System
© 2025 DoctorHealthy1. All rights reserved.`, timestamp, timestamp, time.Now().AddDate(0, 1, 0).Format("January 2, 2006")),
	}

	content, exists := templates[pdfType]
	if !exists {
		return "", fmt.Errorf("unknown PDF type: %s", pdfType)
	}

	return content, nil
}
