from flask import Blueprint, session, jsonify, render_template
from roles import require_role

routes_bp = Blueprint("routes", __name__)

SECTION_DEFINITIONS = [
    {
        "name": "Manage Users",
        "url": "/admin/users",
        "description": "Create, edit, and deactivate employee accounts across the platform.",
        "roles": ["Admin"],
    },
    {
        "name": "View Reports",
        "url": "/reports",
        "description": "Active user counts, weekly signups, and platform health metrics.",
        "roles": ["Admin", "Analyst"],
    },
    {
        "name": "Training Materials",
        "url": "/training",
        "description": "Onboarding guides and nutrition coaching resources for staff.",
        "roles": ["Admin", "Trainer"],
    },
]

ROLE_COLORS = {
    "Admin": "#E14925",     # Orange
    "Analyst": "#3B5BA5",   # slate blue
    "Trainer": "#C9A227",   # gold
    "Employee": "#4C7A5D",  # sage
}

def get_initials(name):
    if not name:
        return "?"
    parts = name.strip().split()
    if len(parts) == 1:
        return parts[0][:2].upper()
    return (parts[0][0] + parts[-1][0]).upper()

@routes_bp.route("/")
def dashboard():
    user = session.get("user")
    if not user:
        return '<a href="/login">Log in with Okta</a>'

    groups = user.get("groups", [])
    all_sections = [
        {**s, "unlocked": any(role in groups for role in s["roles"])}
        for s in SECTION_DEFINITIONS
    ]

    unlocked_count = sum(1 for s in all_sections if s["unlocked"])

    return render_template(
        "dashboard.html",
        user=user,
        groups=groups,
        all_sections=all_sections,
        unlocked_count=unlocked_count,
        total_count=len(all_sections),
        initials=get_initials(user.get("name")),
        role_colors=ROLE_COLORS,
    )
    
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
