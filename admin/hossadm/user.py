import requests
import os


def remove_user(username: str, server_url: str) -> None:
    """Function to call the remove user endpoint

    Args:
        username: username to remove
        server_url: server root endpoint (e.g. https://hoss.mycompany.com)

    Returns:
        None
    """
    pat = os.environ.get("HOSS_PAT")
    if not pat:
        raise Exception("Please set the environment variable HOSS_PAT to a PAT from an admin user.")

    # Find the Auth endpoint
    response = requests.get(f"{server_url}/core/v1/discover")
    if response.status_code != 200:
        raise Exception(f"Failed to reach discover endpoint for server '{server_url}'")

    data = response.json()
    auth_endpoint = data['auth_service']

    # Get a JWT
    headers = {"Authorization": f"Bearer {pat}"}
    try:
        resp = requests.request("POST", f"{auth_endpoint}/pat/exchange/jwt", headers=headers)
    except requests.exceptions.ConnectionError:
        raise Exception(f"Cannot reach Hoss auth service for server '{server_url}'. "
                        f"Verify your network connection and try again.")

    if not resp.ok:
        raise Exception("Could not retrieve JWT using PAT: " + resp.text)

    jwt = resp.json()["id_token"]

    # Remove the user
    headers = {"Authorization": f"Bearer {jwt}"}
    response = requests.delete(f"{auth_endpoint}/user/{username}", headers=headers)
    if response.status_code != 204:
        raise Exception(f"Failed to remove user '{username}' from server '{server_url}'")

    print(f"Successfully removed PATs and group memberships for '{username}'")
