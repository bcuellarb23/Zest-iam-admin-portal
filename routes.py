from flask import Blueprint, session, jsonify

routes_bp = Blueprint("routes", __name__)


@routes_bp.route("/")
def dashboard():
    user = session.get("user")
    if not user:
        return '<a href="/login">Log in with Okta</a>'
    return jsonify({
        "message": f"Welcome to Zest Admin Portal, {user.get('name')}",
        "claims": user
    })


@routes_bp.route("/whoami")
def whoami():
    """Inspect the raw claims Okta sent back, including the groups claim (for RBAC later)."""
    user = session.get("user")
    if not user:
        return jsonify({"error": "not logged in"}), 401
    return jsonify(user)
