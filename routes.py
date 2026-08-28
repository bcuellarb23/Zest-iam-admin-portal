from flask import Blueprint, session, jsonify, render_template, request
from roles import require_role 
import okta_api

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
    try:
        okta_users = okta_api.list_users()
        user_summary = [
            {"email": u["profile"]["email"], "status": u["status"]}
            for u in okta_users
        ]
        return jsonify({
            "message": "User management panel (Admin only)",
            "user_count": len(user_summary),
            "users": user_summary
        })
    except Exception as e:
        return jsonify({"error": "Failed to fetch directory", "detail": str(e)}), 502



@routes_bp.route("/reports")
@require_role("Admin", "Analyst")
def reports():
    user = session.get("user")
    groups = user.get("groups", [])

    event_filter = request.args.get("event_filter", "")
    since = request.args.get("since", "")
    until = request.args.get("until", "")

    try:
        okta_logs = okta_api.list_logs()
        logs_summary = [
            {
                "Event Type": log["eventType"],
                "Date": log.get("published", "N/A"),
                "Actor": log["actor"]["displayName"],
                "Outcome": log["outcome"],
                "Severity": log["severity"],
                "Message": log["displayMessage"],
            }
            for log in okta_logs
        ]

        if event_filter:
            logs_summary = [l for l in logs_summary if l["Event Type"] == event_filter]

        if since:
            logs_summary = [l for l in logs_summary if l["Date"] >= since]

        if until:
            logs_summary = [l for l in logs_summary if l["Date"] <= until + "T23:59:59"]

        event_types = sorted({l["Event Type"] for l in logs_summary})

        return render_template(
            "report.html",
            user=user,
            groups=groups,
            logs=logs_summary,
            event_types=event_types,
            selected_filter=event_filter,
            since=since,
            until=until,
            initials=get_initials(user.get("name")),
            role_colors=ROLE_COLORS,
        )

    except Exception as e:
        return jsonify({"error": "Failed to fetch logs", "detail": str(e)}), 502


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
