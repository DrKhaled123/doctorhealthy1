#!/usr/bin/env python3
"""
GitHub Integration for Factory Orchestrator
Handles automated PR creation, management, and deployment
"""

import os
import json
import time
import asyncio
from typing import Dict, List, Any, Optional
from datetime import datetime
import requests
from urllib.parse import urlparse
import base64

class GitHubIntegration:
    """GitHub integration for automated PR management"""

    def __init__(self, token: str = None, repository: str = None):
        self.token = token or os.getenv('GITHUB_TOKEN')
        self.repository = repository or os.getenv('GITHUB_REPOSITORY')
        self.base_url = "https://api.github.com"
        self.headers = {
            'Authorization': f'token {self.token}',
            'Accept': 'application/vnd.github.v3+json',
            'User-Agent': 'Factory-Orchestrator/1.0'
        }

        if not self.token:
            print("⚠️  GitHub token not provided")
        if not self.repository:
            print("⚠️  GitHub repository not specified")

    def _make_request(self, method: str, endpoint: str, data: Dict = None) -> Dict:
        """Make authenticated request to GitHub API"""
        url = f"{self.base_url}{endpoint}"

        try:
            if method.upper() == 'GET':
                response = requests.get(url, headers=self.headers, params=data)
            elif method.upper() == 'POST':
                response = requests.post(url, headers=self.headers, json=data)
            elif method.upper() == 'PUT':
                response = requests.put(url, headers=self.headers, json=data)
            elif method.upper() == 'PATCH':
                response = requests.patch(url, headers=self.headers, json=data)
            else:
                raise ValueError(f"Unsupported HTTP method: {method}")

            response.raise_for_status()
            return response.json()

        except requests.exceptions.RequestException as e:
            print(f"❌ GitHub API request failed: {e}")
            return {}

    def create_pull_request(self, title: str, body: str, head_branch: str,
                          base_branch: str = "main", labels: List[str] = None) -> Optional[str]:
        """Create a new pull request"""
        if not self.repository or not self.token:
            print("❌ GitHub repository or token not configured")
            return None

        endpoint = f"/repos/{self.repository}/pulls"
        data = {
            "title": title,
            "body": body,
            "head": head_branch,
            "base": base_branch
        }

        if labels:
            data["labels"] = labels

        try:
            response = self._make_request('POST', endpoint, data)

            if response and 'number' in response:
                pr_number = response['number']
                pr_url = response['html_url']
                print(f"✅ Created PR #{pr_number}: {pr_url}")
                return str(pr_number)
            else:
                print("❌ Failed to create PR")
                return None

        except Exception as e:
            print(f"❌ Error creating PR: {e}")
            return None

    def update_pull_request(self, pr_number: str, title: str = None,
                          body: str = None, labels: List[str] = None) -> bool:
        """Update an existing pull request"""
        if not self.repository or not self.token:
            return False

        endpoint = f"/repos/{self.repository}/pulls/{pr_number}"

        data = {}
        if title:
            data["title"] = title
        if body:
            data["body"] = body

        try:
            response = self._make_request('PATCH', endpoint, data)

            if labels:
                # Update labels separately
                label_endpoint = f"/repos/{self.repository}/issues/{pr_number}/labels"
                self._make_request('PUT', label_endpoint, labels)

            print(f"✅ Updated PR #{pr_number}")
            return True

        except Exception as e:
            print(f"❌ Error updating PR #{pr_number}: {e}")
            return False

    def add_pr_comment(self, pr_number: str, comment: str) -> bool:
        """Add a comment to a pull request"""
        if not self.repository or not self.token:
            return False

        endpoint = f"/repos/{self.repository}/issues/{pr_number}/comments"
        data = {"body": comment}

        try:
            response = self._make_request('POST', endpoint, data)
            print(f"✅ Added comment to PR #{pr_number}")
            return True

        except Exception as e:
            print(f"❌ Error adding comment to PR #{pr_number}: {e}")
            return False

    def get_pr_status(self, pr_number: str) -> Dict:
        """Get the status of a pull request"""
        if not self.repository or not self.token:
            return {}

        endpoint = f"/repos/{self.repository}/pulls/{pr_number}"

        try:
            response = self._make_request('GET', endpoint)
            return response

        except Exception as e:
            print(f"❌ Error getting PR status: {e}")
            return {}

    def merge_pull_request(self, pr_number: str, merge_method: str = "merge",
                          commit_message: str = None) -> bool:
        """Merge a pull request"""
        if not self.repository or not self.token:
            return False

        endpoint = f"/repos/{self.repository}/pulls/{pr_number}/merge"
        data = {"merge_method": merge_method}

        if commit_message:
            data["commit_message"] = commit_message

        try:
            response = self._make_request('PUT', endpoint, data)

            if response and response.get('merged'):
                print(f"✅ Merged PR #{pr_number}")
                return True
            else:
                print(f"❌ Failed to merge PR #{pr_number}")
                return False

        except Exception as e:
            print(f"❌ Error merging PR #{pr_number}: {e}")
            return False

    def create_branch(self, branch_name: str, base_branch: str = "main") -> bool:
        """Create a new branch from base branch"""
        if not self.repository or not self.token:
            return False

        # Get the SHA of the base branch
        endpoint = f"/repos/{self.repository}/git/ref/heads/{base_branch}"
        ref_response = self._make_request('GET', endpoint)

        if not ref_response or 'object' not in ref_response:
            print(f"❌ Could not find base branch: {base_branch}")
            return False

        base_sha = ref_response['object']['sha']

        # Create new branch
        create_endpoint = f"/repos/{self.repository}/git/refs"
        branch_data = {
            "ref": f"refs/heads/{branch_name}",
            "sha": base_sha
        }

        try:
            response = self._make_request('POST', create_endpoint, branch_data)
            print(f"✅ Created branch: {branch_name}")
            return True

        except Exception as e:
            print(f"❌ Error creating branch {branch_name}: {e}")
            return False

    def commit_file(self, file_path: str, content: str, branch: str,
                   message: str, author: Dict = None) -> bool:
        """Commit a file to a specific branch"""
        if not self.repository or not self.token:
            return False

        try:
            # Get current file SHA if it exists
            file_endpoint = f"/repos/{self.repository}/contents/{file_path}"
            current_file = self._make_request('GET', f"{file_endpoint}?ref={branch}")

            file_data = {
                "message": message,
                "content": base64.b64encode(content.encode('utf-8')).decode('utf-8'),
                "branch": branch
            }

            if current_file and 'sha' in current_file:
                file_data["sha"] = current_file['sha']

            if author:
                file_data["author"] = author

            response = self._make_request('PUT', file_endpoint, file_data)
            print(f"✅ Committed {file_path} to branch {branch}")
            return True

        except Exception as e:
            print(f"❌ Error committing file {file_path}: {e}")
            return False

    def get_file_content(self, file_path: str, branch: str = "main") -> Optional[str]:
        """Get the content of a file from GitHub"""
        if not self.repository or not self.token:
            return None

        endpoint = f"/repos/{self.repository}/contents/{file_path}?ref={branch}"

        try:
            response = self._make_request('GET', endpoint)

            if response and 'content' in response:
                import base64
                return base64.b64decode(response['content']).decode('utf-8')

        except Exception as e:
            print(f"❌ Error getting file content: {e}")

        return None

    def create_release(self, tag_name: str, release_name: str, body: str,
                      branch: str = "main", prerelease: bool = False) -> Optional[str]:
        """Create a GitHub release"""
        if not self.repository or not self.token:
            return None

        endpoint = f"/repos/{self.repository}/releases"
        data = {
            "tag_name": tag_name,
            "name": release_name,
            "body": body,
            "target_commitish": branch,
            "prerelease": prerelease
        }

        try:
            response = self._make_request('POST', endpoint, data)

            if response and 'id' in response:
                release_id = response['id']
                print(f"✅ Created release: {tag_name}")
                return str(release_id)

        except Exception as e:
            print(f"❌ Error creating release: {e}")

        return None

    def get_workflow_runs(self, branch: str = "main", status: str = None) -> List[Dict]:
        """Get GitHub Actions workflow runs"""
        if not self.repository or not self.token:
            return []

        endpoint = f"/repos/{self.repository}/actions/runs"
        params = {"branch": branch}

        if status:
            params["status"] = status

        try:
            response = self._make_request('GET', endpoint, params)

            if response and 'workflow_runs' in response:
                return response['workflow_runs']

        except Exception as e:
            print(f"❌ Error getting workflow runs: {e}")

        return []

    def trigger_workflow(self, workflow_file: str, branch: str = "main",
                        inputs: Dict = None) -> Optional[str]:
        """Trigger a GitHub Actions workflow"""
        if not self.repository or not self.token:
            return None

        endpoint = f"/repos/{self.repository}/actions/workflows/{workflow_file}/dispatches"
        data = {
            "ref": branch
        }

        if inputs:
            data["inputs"] = inputs

        try:
            response = self._make_request('POST', endpoint, data)

            if response:
                print(f"✅ Triggered workflow: {workflow_file}")
                return "triggered"

        except Exception as e:
            print(f"❌ Error triggering workflow: {e}")

        return None

# Factory integration functions
def create_factory_pr(factory_result: Dict, branch_name: str = None) -> Optional[str]:
    """Create a PR from factory results"""
    github = GitHubIntegration()

    if not branch_name:
        branch_name = f"factory-{int(time.time())}"

    # Create branch
    if not github.create_branch(branch_name):
        return None

    # Commit generated files
    if 'generated_files' in factory_result:
        for file_path, content in factory_result['generated_files'].items():
            github.commit_file(
                file_path=file_path,
                content=content,
                branch=branch_name,
                message=f"Factory: Add {file_path}"
            )

    # Create PR
    pr_title = factory_result.get('pr_title', 'Factory Generated Changes')
    pr_body = factory_result.get('pr_description', 'Automated changes from Factory Orchestrator')

    pr_number = github.create_pull_request(
        title=pr_title,
        body=pr_body,
        head_branch=branch_name,
        labels=['factory', 'automated']
    )

    return pr_number

def update_pr_with_test_results(pr_number: str, test_results: Dict) -> bool:
    """Update PR with test results"""
    github = GitHubIntegration()

    # Format test results as comment
    comment = "## 🧪 Test Results\n\n"
    comment += f"**Status**: {'✅ PASSED' if test_results['success'] else '❌ FAILED'}\n"
    comment += f"**Tests Run**: {test_results.get('tests_run', 0)}\n"
    comment += f"**Passed**: {test_results.get('passed', 0)}\n"
    comment += f"**Failed**: {test_results.get('failed', 0)}\n"
    comment += f"**Duration**: {test_results.get('duration', 0):.2f}s\n\n"

    if 'coverage' in test_results:
        comment += f"**Coverage**: {test_results['coverage']".1%"}\n\n"

    if 'details' in test_results:
        comment += "**Details**:\n"
        for detail in test_results['details'][:10]:  # Limit to 10 details
            comment += f"- {detail}\n"

    return github.add_pr_comment(pr_number, comment)

def deploy_from_pr(pr_number: str, environment: str = "production") -> bool:
    """Trigger deployment from PR"""
    github = GitHubIntegration()

    # Get PR details
    pr_details = github.get_pr_status(pr_number)

    if not pr_details:
        return False

    # Trigger deployment workflow
    workflow_result = github.trigger_workflow(
        workflow_file="deploy.yml",
        branch=pr_details.get('head', {}).get('ref', 'main'),
        inputs={
            "pr_number": str(pr_number),
            "environment": environment,
            "auto_deploy": "true"
        }
    )

    return workflow_result is not None

# Example usage
async def demo_github_integration():
    """Demonstrate GitHub integration"""
    print("🐙 GitHub Integration Demo")
    print("=" * 40)

    # Initialize GitHub integration
    github = GitHubIntegration()

    if not github.token or not github.repository:
        print("❌ GitHub token or repository not configured")
        print("Set GITHUB_TOKEN and GITHUB_REPOSITORY environment variables")
        return

    # Example: Create a PR from factory results
    factory_result = {
        'pr_title': 'Factory: Add user authentication system',
        'pr_description': '''## 🤖 Factory Generated Changes

This PR contains automated changes from the Factory Orchestrator:

### Changes Made:
- Added user authentication system
- Implemented JWT token handling
- Added password validation
- Created login/logout endpoints

### Testing:
- Unit tests: ✅ 15/15 passed
- Integration tests: ✅ 8/8 passed
- Security scan: ✅ No vulnerabilities found

### Auto-generated by Factory Orchestrator''',
        'generated_files': {
            'auth/user_auth.py': '''def authenticate_user(username: str, password: str) -> bool:
    # Authentication logic
    return True''',
            'auth/test_auth.py': '''def test_authentication():
    # Test logic
    assert True'''
        }
    }

    # Create PR
    pr_number = create_factory_pr(factory_result, "feature/user-auth")

    if pr_number:
        print(f"✅ Created PR #{pr_number}")

        # Add test results
        test_results = {
            'success': True,
            'tests_run': 23,
            'passed': 23,
            'failed': 0,
            'duration': 45.2,
            'coverage': 92.5,
            'details': ['All authentication tests passed', 'Security scan clean']
        }

        if update_pr_with_test_results(pr_number, test_results):
            print("✅ Added test results to PR")

        # Trigger deployment
        if deploy_from_pr(pr_number, "staging"):
            print("✅ Triggered deployment from PR")

    print("\n✅ GitHub integration demo completed!")

if __name__ == "__main__":
    asyncio.run(demo_github_integration())