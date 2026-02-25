import { Outlet } from "react-router-dom";
import { Activity } from "lucide-react";
import { motion } from "framer-motion";

export default function AuthLayout() {
    return (
        <div className="flex min-h-screen bg-background text-foreground">
            {/* Left Pane - Branding & Animation */}
            <div className="hidden lg:flex flex-col flex-1 bg-muted/20 relative overflow-hidden items-center justify-center p-12 border-r border-border/50">
                <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent pointer-events-none" />
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.8 }}
                    className="relative z-10 text-center"
                >
                    <div className="flex items-center justify-center gap-3 mb-6">
                        <Activity className="h-12 w-12 text-primary" />
                        <h1 className="text-4xl font-bold tracking-tight">SentinelAI</h1>
                    </div>
                    <p className="text-lg text-muted-foreground max-w-md mx-auto">
                        Production-grade distributed target monitoring layered with AI failure analysis orchestration.
                    </p>
                </motion.div>
            </div>

            {/* Right Pane - Auth Form */}
            <div className="flex-1 flex items-center justify-center p-6 relative">
                <div className="absolute top-0 w-full h-1 bg-gradient-to-r from-transparent via-primary/50 to-transparent" />
                <div className="w-full max-w-md">
                    <Outlet />
                </div>
            </div>
        </div>
    );
}
