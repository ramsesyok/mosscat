import os
import subprocess
import time
from pathlib import Path

import pytest
import requests

REPO_ROOT = Path(__file__).resolve().parents[1]
BINARY = REPO_ROOT / "mosscat-test.exe"

@pytest.fixture(scope="session", autouse=True)
def server():
    # remove old database file if exists
    db_file = REPO_ROOT / "mosscat.db"
    if db_file.exists():
        db_file.unlink()
    
    # build binary
    subprocess.run(["go", "build", "-o", str(BINARY), "."], cwd=REPO_ROOT, check=True)
    proc = subprocess.Popen([str(BINARY)], cwd=REPO_ROOT)
    passwd_file = REPO_ROOT / "admin.initial.password"
    # wait for password file
    for _ in range(30):
        if passwd_file.exists():
            break
        time.sleep(1)
    else:
        proc.terminate()
        pytest.fail("server did not start")
    password = passwd_file.read_text().strip()
    base_url = "http://127.0.0.1:8080"
    os.environ["BASE_URL"] = f"{base_url}/api"
    os.environ["ADMIN_PASSWORD"] = password
    yield base_url, password
    proc.terminate()
    proc.wait()
    if BINARY.exists():
        BINARY.unlink()
    if passwd_file.exists():
        passwd_file.unlink()
    if db_file.exists():
        db_file.unlink()

@pytest.fixture(scope="session", autouse=True)
def token(server):
    base_url, password = server
    res = requests.post(f"{base_url}/api/auth/login", json={"username": "admin", "password": password})
    res.raise_for_status()
    tok = res.json()["accessToken"]
    os.environ["TOKEN"] = tok
    return tok

@pytest.fixture(scope="session")
def base_url(server):
    return server[0]

@pytest.fixture(scope="session")
def admin_password(server):
    return server[1]
