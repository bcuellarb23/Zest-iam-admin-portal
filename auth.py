import os
from flask import Blueprint, redirect, url_for, session
from authlib.integrations.flask_client import OAuth

auth_bp = Blueprint("auth", __name__)

oauth = OAuth()


def init_oauth(app):
    """Call this once from app.py to wire OAuth into this blueprint."""
    oauth.init_app(app)
    oauth.register(
        name="okta",
        client_id=os.getenv("OKTA_CLIENT_ID"),
        client_secret=os.getenv("OKTA_CLIENT_SECRET"),
        server_metadata_url=f"{os.getenv('OKTA_ISSUER')}/.well-known/openid-configuration",
        client_kwargs={"scope": "openid profile email"},
    )


@auth_bp.route("/login")
def login():
    redirect_uri = url_for("auth.callback", _external=True)
    return oauth.okta.authorize_redirect(redirect_uri)


# NOTE: path matches Okta's default redirect URI for this app:
# http://localhost:8080/authorization-code/callback
@auth_bp.route("/authorization-code/callback")
def callback():
    token = oauth.okta.authorize_access_token()
    user_info = token.get("userinfo")
    session["user"] = user_info
    return redirect(url_for("routes.dashboard"))


@auth_bp.route("/logout")
def logout():
    session.pop("user", None)
    issuer = os.getenv("OKTA_ISSUER")
    return redirect(
        f"{issuer}/v1/logout?id_token_hint=&post_logout_redirect_uri={url_for('routes.dashboard', _external=True)}"
    )
