package handlers

import "github.com/gofiber/fiber/v2"

// Dashboard routes
func DashboardHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f5f5f5; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 1rem 2rem; }
        .header h1 { font-size: 1.8rem; margin-bottom: 0.5rem; }
        .container { padding: 2rem; max-width: 1200px; margin: 0 auto; }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1.5rem; margin-bottom: 2rem; }
        .stat-card { background: white; padding: 1.5rem; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .stat-card h3 { color: #333; margin-bottom: 1rem; font-size: 1rem; }
        .stat-number { font-size: 2rem; font-weight: bold; color: #667eea; }
        .section { background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-bottom: 2rem; }
        .section-header { padding: 1.5rem; border-bottom: 1px solid #eee; }
        .section-header h2 { color: #333; }
        .section-content { padding: 1.5rem; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 0.75rem; text-align: left; border-bottom: 1px solid #eee; }
        th { background-color: #f8f9fa; font-weight: 600; }
        .status-active { color: #28a745; font-weight: 600; }
        .status-inactive { color: #dc3545; font-weight: 600; }
        .btn { padding: 0.5rem 1rem; border: none; border-radius: 5px; cursor: pointer; font-size: 0.9rem; }
        .btn-primary { background: #667eea; color: white; }
        .btn-danger { background: #dc3545; color: white; }
        .btn:hover { opacity: 0.8; }
        .refresh-btn { float: right; margin-top: -0.5rem; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🚀 Admin Dashboard</h1>
        <p>Monitor and manage your application</p>
    </div>
    
    <div class="container">
        <div class="stats-grid" id="stats-grid">
            <!-- Stats will be loaded here -->
        </div>
        
        <div class="section">
            <div class="section-header">
                <h2>👥 User Management</h2>
                <button class="btn btn-primary refresh-btn" onclick="loadUsers()">Refresh</button>
            </div>
            <div class="section-content">
                <table id="users-table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Username</th>
                            <th>Email</th>
                            <th>Role</th>
                            <th>Status</th>
                            <th>Last Login</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody id="users-tbody">
                        <!-- Users will be loaded here -->
                    </tbody>
                </table>
            </div>
        </div>
        
        <div class="section">
            <div class="section-header">
                <h2>📊 Recent Activity</h2>
                <button class="btn btn-primary refresh-btn" onclick="loadActivity()">Refresh</button>
            </div>
            <div class="section-content">
                <table id="activity-table">
                    <thead>
                        <tr>
                            <th>Time</th>
                            <th>Action</th>
                            <th>User</th>
                            <th>Details</th>
                        </tr>
                    </thead>
                    <tbody id="activity-tbody">
                        <!-- Activity will be loaded here -->
                    </tbody>
                </table>
            </div>
        </div>
    </div>
    
    <script>
        // Load dashboard data
        async function loadStats() {
            try {
                const response = await fetch('/api/stats');
                const stats = await response.json();
                
                document.getElementById('stats-grid').innerHTML = ` + "`" + `
                    <div class="stat-card">
                        <h3>👥 Total Users</h3>
                        <div class="stat-number">${stats.total_users}</div>
                    </div>
                    <div class="stat-card">
                        <h3>✅ Active Users</h3>
                        <div class="stat-number">${stats.active_users}</div>
                    </div>
                    <div class="stat-card">
                        <h3>❌ Inactive Users</h3>
                        <div class="stat-number">${stats.inactive_users}</div>
                    </div>
                    <div class="stat-card">
                        <h3>💾 Memory Usage</h3>
                        <div class="stat-number">${stats.memory_usage}</div>
                    </div>
                    <div class="stat-card">
                        <h3>⏱️ Uptime</h3>
                        <div class="stat-number" style="font-size: 1.2rem;">${stats.uptime}</div>
                    </div>
                ` + "`" + `;
            } catch (error) {
                console.error('Error loading stats:', error);
            }
        }
        
        async function loadUsers() {
            try {
                const response = await fetch('/api/users');
                const users = await response.json();
                
                const tbody = document.getElementById('users-tbody');
                tbody.innerHTML = users.map(user => ` + "`" + `
                    <tr>
                        <td>${user.id}</td>
                        <td>${user.username}</td>
                        <td>${user.email}</td>
                        <td>${user.role}</td>
                        <td><span class="status-${user.status.toLowerCase()}">${user.status}</span></td>
                        <td>${new Date(user.last_login).toLocaleString()}</td>
                        <td>
                            <button class="btn btn-primary" onclick="toggleUserStatus(${user.id})">
                                ${user.status === 'Active' ? 'Deactivate' : 'Activate'}
                            </button>
                        </td>
                    </tr>
                ` + "`" + `).join('');
            } catch (error) {
                console.error('Error loading users:', error);
            }
        }
        
        async function loadActivity() {
            try {
                const response = await fetch('/api/activity');
                const activity = await response.json();
                
                const tbody = document.getElementById('activity-tbody');
                tbody.innerHTML = activity.slice(0, 10).map(log => ` + "`" + `
                    <tr>
                        <td>${new Date(log.timestamp).toLocaleString()}</td>
                        <td>${log.action}</td>
                        <td>${log.username}</td>
                        <td>${log.details}</td>
                    </tr>
                ` + "`" + `).join('');
            } catch (error) {
                console.error('Error loading activity:', error);
            }
        }
        
        async function toggleUserStatus(userId) {
            try {
                const response = await fetch(` + "`" + `/api/users/${userId}/toggle` + "`" + `, {
                    method: 'POST'
                });
                if (response.ok) {
                    loadUsers();
                    loadStats();
                    loadActivity();
                }
            } catch (error) {
                console.error('Error toggling user status:', error);
            }
        }
        
        // Load data on page load
        document.addEventListener('DOMContentLoaded', function() {
            loadStats();
            loadUsers();
            loadActivity();
            
            // Auto-refresh every 30 seconds
            setInterval(() => {
                loadStats();
            }, 30000);
        });
    </script>
</body>
</html>
	`)
}
