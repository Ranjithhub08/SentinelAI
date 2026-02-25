import { Outlet, Navigate, Link, useLocation } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { LogOut, LayoutDashboard, Activity, Settings } from "lucide-react";

export default function DashboardLayout() {
    const { isAuthenticated, logout } = useAuth();
    const location = useLocation();

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    const navItems = [
        { name: "Dashboard", href: "/", icon: LayoutDashboard },
        { name: "Monitors", href: "/", icon: Activity },
        { name: "Settings", href: "/", icon: Settings },
    ];

    return (
        <div className="min-h-screen flex bg-background text-foreground">
            {/* Sidebar Navigation */}
            <aside className="w-64 border-r border-border/50 bg-card/30 flex flex-col">
                <div className="h-16 flex items-center px-6 border-b border-border/50">
                    <Activity className="h-6 w-6 text-primary mr-2" />
                    <span className="font-semibold text-lg tracking-tight">SentinelAI</span>
                </div>
                <nav className="flex-1 py-6 px-4 space-y-2">
                    {navItems.map((item) => {
                        const isActive = location.pathname === item.href;
                        return (
                            <Link
                                key={item.name}
                                to={item.href}
                                className={`flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors ${isActive
                                        ? "bg-primary/10 text-primary font-medium"
                                        : "text-muted-foreground hover:bg-muted/50 hover:text-foreground"
                                    }`}
                            >
                                <item.icon className="h-4 w-4" />
                                {item.name}
                            </Link>
                        );
                    })}
                </nav>
                <div className="p-4 border-t border-border/50">
                    <button
                        onClick={logout}
                        className="flex items-center gap-3 w-full px-3 py-2 rounded-md text-sm text-muted-foreground hover:bg-destructive/10 hover:text-destructive transition-colors"
                    >
                        <LogOut className="h-4 w-4" />
                        Sign Out
                    </button>
                </div>
            </aside>

            {/* Main Content Area */}
            <main className="flex-1 flex flex-col min-h-screen overflow-auto">
                <header className="h-16 flex items-center justify-between px-8 border-b border-border/50 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-10">
                    <h2 className="text-lg font-medium">Workspace</h2>
                    <div className="h-8 w-8 rounded-full bg-primary/20 flex items-center justify-center text-primary text-sm font-medium">
                        EN
                    </div>
                </header>
                <div className="p-8 flex-1">
                    <Outlet />
                </div>
            </main>
        </div>
    );
}
