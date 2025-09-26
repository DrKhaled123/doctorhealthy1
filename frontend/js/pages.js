async function apiPost(url, apiKey, payload) {
  const res = await fetch(url + (url.includes('?') ? '&' : '?') + 'language=' + encodeURIComponent(payload.language || 'en'), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': apiKey || ''
    },
    body: JSON.stringify(payload)
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`HTTP ${res.status}: ${text}`);
  }
  return res.json();
}

async function apiPostPdf(url, apiKey, payload) {
  const res = await fetch(url + (url.includes('?') ? '&' : '?') + 'language=' + encodeURIComponent(payload.language || 'en'), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': apiKey || ''
    },
    body: JSON.stringify(payload)
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`HTTP ${res.status}: ${text}`);
  }
  return res.blob();
}

function downloadBlob(blob, filename) {
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url; a.download = filename; a.click();
  setTimeout(() => window.URL.revokeObjectURL(url), 1000);
}

function wireCommon(formId) {
  const form = document.getElementById(formId);
  const apiKeyInput = document.getElementById('apiKey');
  const result = document.getElementById('result');
  const pre = document.getElementById('resultPre');
  const pdfBtn = document.getElementById('downloadPdf');
  return { form, apiKeyInput, result, pre, pdfBtn };
}

function parseForm(form) {
  const data = new FormData(form);
  const obj = {};
  data.forEach((v,k) => { obj[k] = typeof v === 'string' ? v.trim() : v; });
  // minimal normalization
  if (obj.symptoms) obj.symptoms = obj.symptoms.split(',').map(s => s.trim()).filter(Boolean);
  if (obj.dislikes) obj.dislikes = obj.dislikes.split(',').map(s => s.trim()).filter(Boolean);
  if (obj.height_cm) obj.height_cm = Number(obj.height_cm);
  if (obj.weight_kg) obj.weight_kg = Number(obj.weight_kg);
  if (obj.age) obj.age = Number(obj.age);
  if (obj.sessions) obj.sessions = Number(obj.sessions);
  return obj;
}

function setupDietPage(cfg) {
  const { form, apiKeyInput, result, pre, pdfBtn } = wireCommon('dietForm');
  let lastPayload = null;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const payload = parseForm(form); lastPayload = payload;
    try {
      const json = await apiPost(cfg.generateUrl, apiKeyInput.value, payload);
      pre.textContent = JSON.stringify(json, null, 2);
      result.hidden = false; pdfBtn.disabled = false;
    } catch (err) {
      pre.textContent = String(err); result.hidden = false; pdfBtn.disabled = true;
    }
  });
  pdfBtn.addEventListener('click', async () => {
    if (!lastPayload) return;
    try {
      const blob = await apiPostPdf(cfg.pdfUrl, apiKeyInput.value, lastPayload);
      downloadBlob(blob, 'diet_plan.pdf');
    } catch (err) {
      alert(String(err));
    }
  });
}

function setupWorkoutPage(cfg) {
  const { form, apiKeyInput, result, pre, pdfBtn } = wireCommon('workoutForm');
  let lastPayload = null;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const payload = parseForm(form); lastPayload = payload;
    try {
      const json = await apiPost(cfg.generateUrl, apiKeyInput.value, payload);
      pre.textContent = JSON.stringify(json, null, 2);
      result.hidden = false; pdfBtn.disabled = false;
    } catch (err) {
      pre.textContent = String(err); result.hidden = false; pdfBtn.disabled = true;
    }
  });
  pdfBtn.addEventListener('click', async () => {
    if (!lastPayload) return;
    try {
      const blob = await apiPostPdf(cfg.pdfUrl, apiKeyInput.value, lastPayload);
      downloadBlob(blob, 'workout_plan.pdf');
    } catch (err) {
      alert(String(err));
    }
  });
}

function setupLifestylePage(cfg) {
  const { form, apiKeyInput, result, pre, pdfBtn } = wireCommon('lifestyleForm');
  let lastPayload = null;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const payload = parseForm(form); lastPayload = payload;
    try {
      const json = await apiPost(cfg.generateUrl, apiKeyInput.value, payload);
      pre.textContent = JSON.stringify(json, null, 2);
      result.hidden = false; pdfBtn.disabled = false;
    } catch (err) {
      pre.textContent = String(err); result.hidden = false; pdfBtn.disabled = true;
    }
  });
  pdfBtn.addEventListener('click', async () => {
    if (!lastPayload) return;
    try {
      const blob = await apiPostPdf(cfg.pdfUrl, apiKeyInput.value, lastPayload);
      downloadBlob(blob, 'lifestyle_plan.pdf');
    } catch (err) {
      alert(String(err));
    }
  });
}

function setupRecipesPage(cfg) {
  const { form, apiKeyInput, result, pre, pdfBtn } = wireCommon('recipesForm');
  let lastPayload = null;
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    const payload = parseForm(form); lastPayload = payload;
    try {
      const json = await apiPost(cfg.generateUrl, apiKeyInput.value, payload);
      pre.textContent = JSON.stringify(json, null, 2);
      result.hidden = false; pdfBtn.disabled = false;
    } catch (err) {
      pre.textContent = String(err); result.hidden = false; pdfBtn.disabled = true;
    }
  });
  pdfBtn.addEventListener('click', async () => {
    if (!lastPayload) return;
    try {
      const blob = await apiPostPdf(cfg.pdfUrl, apiKeyInput.value, lastPayload);
      downloadBlob(blob, 'recipes.pdf');
    } catch (err) {
      alert(String(err));
    }
  });
}

// Expose setup functions
window.setupDietPage = setupDietPage;
window.setupWorkoutPage = setupWorkoutPage;
window.setupLifestylePage = setupLifestylePage;
window.setupRecipesPage = setupRecipesPage;


