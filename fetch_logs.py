import urllib.request
import json

url = "https://api.github.com/repos/vysortech/evo-crm-vysor/actions/runs?per_page=1"
try:
    req = urllib.request.Request(url)
    with urllib.request.urlopen(req) as response:
        data = json.loads(response.read().decode())
        if data["workflow_runs"]:
            run_id = data["workflow_runs"][0]["id"]
            print(f"Latest run ID: {run_id}")
            
            jobs_url = f"https://api.github.com/repos/vysortech/evo-crm-vysor/actions/runs/{run_id}/jobs"
            jobs_req = urllib.request.Request(jobs_url)
            with urllib.request.urlopen(jobs_req) as jobs_response:
                jobs_data = json.loads(jobs_response.read().decode())
                for job in jobs_data["jobs"]:
                    print(f"Job {job['name']} status: {job['status']} conclusion: {job['conclusion']}")
                    for step in job["steps"]:
                        if step["conclusion"] == "failure":
                            print(f"  Step '{step['name']}' failed!")
        else:
            print("No workflow runs found.")
except urllib.error.HTTPError as e:
    print(f"HTTP Error: {e.code} - {e.read().decode()}")
except Exception as e:
    print(f"Error: {e}")
