import subprocess
import json

def run_service(bucket, object_name):
    env = {
        "GCP_BUCKET": bucket,
        "GCP_OBJECT": object_name
    }
    result = subprocess.run(["./golden_service"], capture_output=True, text=True, env=env)
    if result.returncode != 0:
        raise RuntimeError(f"Service failed: {result.stderr}")
    return json.loads(result.stdout)

if __name__ == "__main__":
    bucket_name = "your-bucket"
    object_name = "your-object.json"
    output = run_service(bucket_name, object_name)
    print(json.dumps(output, indent=2))
