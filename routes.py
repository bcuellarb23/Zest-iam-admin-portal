from flask import Blueprint, session, jsonify
from roles import require_role

routes_bp = Blueprint("routes", __name__)


@routes_bp.route("/")
def dashboard():
    user = session.get("user")
    if not user:
        return '<a href="/login">Log in with Okta</a>'

    groups = user.get("groups", [])
    available_sections = []

    if "Admin" in groups:
        available_sections.append({"name": "Manage Users", "url": "/admin/users"})
    if "Analyst" in groups:
        available_sections.append({"name": "View Reports", "url": "/reports"})
    if "Trainer" in groups:
        available_sections.append({"name": "Training Materials", "url": "/training"})

    return jsonify({
        "message": f"Welcome to Zest Admin Portal, {user.get('name')}",
        "your_groups": groups,
        "available_sections": available_sections
    })


@routes_bp.route("/admin/users")
@require_role("Admin")
def manage_users():
    return jsonify({
        "message": "User management panel (Admin only)",
        "sample_data": ["user1@zest.com", "user2@zest.com", "user3@zest.com"]
    })


@routes_bp.route("/reports")
@require_role("Admin", "Analyst")
def reports():
    return jsonify({
        "message": "Reports dashboard (Admin or Analyst)",
        "sample_data": {"active_users": 1240, "signups_this_week": 87}
    })


@routes_bp.route("/training")
@require_role("Trainer", "Admin")
def training():
    return jsonify({
        "message": "Training materials (Trainer or Admin)",
        "sample_data": ["Onboarding Guide", "Nutrition Basics Course"]
    })


@routes_bp.route("/whoami")
def whoami():
    """Inspect the raw claims Okta sent back, including the groups claim (for RBAC later)."""
    user = session.get("user")
    if not user:
        return jsonify({"error": "not logged in"}), 401
    return jsonify(user)
