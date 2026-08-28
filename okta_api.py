import os  
import time 
import uuid 
import jwt 
import requests 

OKTA_DOMAIN = os.getenv("OKTA_DOMAIN")
SERVICE_CLIENT_ID = os.getenv("OKTA_SERVICE_CLIENT_ID")
PRIVATE_KEY_PATH = os.getenv("OKTA_PRIVATE_KEY_PATH", "keys/private.pem")
KEY_ID = os.getenv("OKTA_KEY_ID")

TOKEN_URL = f"https://{OKTA_DOMAIN}/oauth2/v1/token"

def _load_private_key():
    with open(PRIVATE_KEY_PATH, "r") as f:
        return f.read()

def _build_client_assertion():
    """
        Builds and sign the JWT that proves this app identity to Okta. 
        This is nor an acess token.
    """

    now = int(time.time())
    
    claims = {
        "iss": SERVICE_CLIENT_ID,
        "sub": SERVICE_CLIENT_ID,
        "aud": TOKEN_URL,
        "exp": now + 300,
        "jti": str(uuid.uuid4()),
    }
    private_key = _load_private_key()
    return jwt.encode(
        claims,
        private_key,
        algorithm="RS256",
        headers={"kid": KEY_ID}
    )

def get_access_token(scope="okta.users.read"):
    """
        Exchange a freshly-signed client assertion for a real Okta access token.
        This have to be called everytime is needes since assertion is a single-use
    """
    
    client_assertion = _build_client_assertion()

    token_request_body = {
        "grant_type": "client_credentials",
        "scope": scope,
        "client_assertion_type": "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
        "client_assertion": client_assertion,
    }
    headers = {
        "Accept": "application/json",
        "Content-Type": "application/x-www-form-urlencoded"
    }
    
    response = requests.post(TOKEN_URL, headers=headers, data=token_request_body)
    # print(response.status_code, response.text)
    response.raise_for_status()
    return response.json()["access_token"]

def list_users():   
    """
        This function calls Okta's Users API usig a fresh obtained access token.
    """

    access_token = get_access_token(scope="okta.users.read")

    url = f"https://{OKTA_DOMAIN}/api/v1/users"
    headers = {
        "Authorization": f"Bearer {access_token}",
        "Accept": "application/json"
    }

    response = requests.get(url, headers=headers)
    """
        This is a debug line
        print(response.status_code, response.text)
    """
    response.raise_for_status()
    return response.json()

def list_logs():
    """
        This function calls Okta's API to get and return all the logs
    """

    access_token = get_access_token(scope="okta.logs.read")

    url = f"https://{OKTA_DOMAIN}/api/v1/logs"
    headers = {
        "Authorization": f"Bearer {access_token}",
        "Accept": "application/json"
    }
    query = {
        "filter": 'actor.displayName co "Bryan"',
        "limit": "50",
        "sortOrder": "DESCENDING"
    }

    response = requests.get(url, headers=headers, params=query)
    try:
        response.raise_for_status()
    except requests.exceptions.HTTPerror as e:
        raise 
    
    return response.json()

