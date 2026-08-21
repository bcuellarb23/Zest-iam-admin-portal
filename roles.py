from functools import wraps
from flask import session, jsonify, render_template, request, redirect, url_for

def require_role(*allowed_roles):
    """
        Restricts a route to users whose Okta groups include at least one
        of the allowed_roles. Usage:

            @app.route("/admin")
            @require_role("Admin")
            def admin_page():
                ...

            @app.route("/reports")
            @require_role("Admin", "Analyst")   # either role can access
            def reports():
                ...
    """
    def decorator(f):
        @wraps(f)
        def wrapped(*args, **kwargs):
            user = session.get("user")

            if not user:
                return jsonify({"error": "not logged in"}), 401

            user_groups = user.get("groups", [])

            if not any(role in user_groups for role in allowed_roles):
                if request.accept_mimetypes.accept_html:
                    return render_template(
                        "forbidden.html",
                        required_roles=",".join(allowed_roles),
                        user_roles=",".join(user_groups) if user_groups else "None"
                ), 403

            return f(*args, **kwargs)
        return wrapped
    return decorator

