import { useState, useEffect } from "react";
import { Plus, Activity as ActivityIcon, Brain, ChevronDown } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { formatDistanceToNow } from "date-fns";
import api from "@/services/api";
import { Monitor } from "@/types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";

export default function Dashboard() {
    const [monitors, setMonitors] = useState<Monitor[]>([]);
    const [isAdding, setIsAdding] = useState(false);
    const [newUrl, setNewUrl] = useState("");
    const [newInterval, setNewInterval] = useState("60");
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchMonitors = async () => {
        try {
            setError(null);
            const res = await api.get("/monitor/list");
            if (res.data.success) {
                setMonitors(res.data.data || []);
            }
        } catch (err) {
            setError("Unable to connect to telemetry stream.");
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchMonitors();
        const interval = setInterval(fetchMonitors, 5000); // 5 sec auto refresh
        return () => clearInterval(interval);
    }, []);

    const handleAddMonitor = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const res = await api.post("/monitor/add", {
                url: newUrl,
                interval: parseInt(newInterval, 10),
            });
            if (res.data.success) {
                setIsAdding(false);
                setNewUrl("");
                fetchMonitors();
            }
        } catch (err) {
            setError("Failed to provision new target.");
        }
    };

    return (
        <div className="max-w-6xl mx-auto space-y-8">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Active Monitors</h1>
                    <p className="text-muted-foreground mt-1">Real-time infrastructure pulse tracked effectively.</p>
                </div>
                <Button onClick={() => setIsAdding(!isAdding)} className="gap-2">
                    <Plus className="h-4 w-4" />
                    Provision Monitor
                </Button>
            </div>

            {error && (
                <div className="p-4 rounded-md border border-destructive/20 bg-destructive/10 text-sm text-destructive-foreground">
                    {error}
                </div>
            )}

            <AnimatePresence>
                {isAdding && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: "auto" }}
                        exit={{ opacity: 0, height: 0 }}
                        className="overflow-hidden"
                    >
                        <Card className="border-border/50 bg-muted/10 mb-8 w-full max-w-xl">
                            <CardContent className="pt-6">
                                <form onSubmit={handleAddMonitor} className="space-y-4">
                                    <div className="grid grid-cols-4 gap-4">
                                        <div className="col-span-3 space-y-2">
                                            <Label htmlFor="url">Target URL</Label>
                                            <Input id="url" type="url" placeholder="https://api.example.com" value={newUrl} onChange={(e) => setNewUrl(e.target.value)} required />
                                        </div>
                                        <div className="col-span-1 space-y-2">
                                            <Label htmlFor="interval">Interval (s)</Label>
                                            <Input id="interval" type="number" min="5" value={newInterval} onChange={(e) => setNewInterval(e.target.value)} required />
                                        </div>
                                    </div>
                                    <div className="flex justify-end gap-2">
                                        <Button type="button" onClick={() => setIsAdding(false)} className="bg-secondary text-secondary-foreground hover:bg-secondary/80 px-4 py-2 rounded-md text-sm">Cancel</Button>
                                        <Button type="submit">Inject Target</Button>
                                    </div>
                                </form>
                            </CardContent>
                        </Card>
                    </motion.div>
                )}
            </AnimatePresence>

            {isLoading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {[1, 2, 3].map((i) => <div key={i} className="h-48 rounded-lg bg-muted/20 animate-pulse" />)}
                </div>
            ) : monitors.length === 0 ? (
                <div className="text-center py-24 border border-dashed rounded-lg border-border/50 bg-muted/5">
                    <ActivityIcon className="mx-auto h-12 w-12 text-muted-foreground/50 mb-4" />
                    <h3 className="text-lg font-medium text-foreground">No target definitions deployed</h3>
                    <p className="text-sm text-muted-foreground mt-1">Start by provisioning a telemetry tracker securely.</p>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                    {monitors.map((m) => (
                        <MonitorCard key={m.id} monitor={m} />
                    ))}
                </div>
            )}
        </div>
    );
}

function MonitorCard({ monitor }: { monitor: Monitor }) {
    const [expanded, setExpanded] = useState(false);
    const isHealthy = monitor.is_healthy;

    return (
        <motion.div layout initial={{ opacity: 0, scale: 0.98 }} animate={{ opacity: 1, scale: 1 }}>
            <Card className={`overflow-hidden transition-all duration-300 ${!isHealthy && monitor.status_code !== 0 ? 'border-destructive/50 shadow-[0_0_15px_rgba(220,38,38,0.1)]' : 'border-border/50 hover:border-primary/30'}`}>
                <div className="p-5">
                    <div className="flex items-start justify-between">
                        <div className="space-y-1 overflow-hidden">
                            <h3 className="font-semibold text-base truncate pr-4" title={monitor.url}>{monitor.url}</h3>
                            <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                <span className="flex items-center gap-1">
                                    <div className={`w-2 h-2 rounded-full ${isHealthy ? 'bg-emerald-500 animate-pulse-slow' : 'bg-destructive animate-pulse-fast'}`} />
                                    {isHealthy ? 'Operational' : 'Failed'}
                                </span>
                                <span>•</span>
                                <span>{monitor.status_code || '--'}</span>
                                <span>•</span>
                                <span>{monitor.response_time ? `${monitor.response_time}ms` : '-- ms'}</span>
                            </div>
                        </div>
                    </div>
                    <div className="mt-4 flex items-center justify-between text-xs text-muted-foreground">
                        <span>Interval: {monitor.interval / 1000000000}s</span>
                        <span>{monitor.last_checked ? formatDistanceToNow(new Date(monitor.last_checked), { addSuffix: true }) : 'Pending'}</span>
                    </div>
                </div>

                {/* AI Insight Dropdown mapped over error context directly */}
                {monitor.ai_explanation && (
                    <div className="border-t border-border/50 bg-muted/10">
                        <button
                            onClick={() => setExpanded(!expanded)}
                            className="w-full text-left px-5 py-3 text-xs font-medium text-muted-foreground hover:text-foreground hover:bg-muted/20 transition-colors flex justify-between items-center"
                        >
                            <span className="flex items-center gap-2 text-primary/80">
                                <Brain className="h-3.5 w-3.5" />
                                AI Analysis Available
                            </span>
                            <ChevronDown className={`h-4 w-4 transition-transform ${expanded ? "rotate-180" : ""}`} />
                        </button>
                        <AnimatePresence>
                            {expanded && (
                                <motion.div
                                    initial={{ height: 0, opacity: 0 }}
                                    animate={{ height: "auto", opacity: 1 }}
                                    exit={{ height: 0, opacity: 0 }}
                                    className="overflow-hidden"
                                >
                                    <div className="p-5 pt-0 text-sm leading-relaxed text-muted-foreground whitespace-pre-wrap">
                                        {monitor.ai_explanation}
                                    </div>
                                </motion.div>
                            )}
                        </AnimatePresence>
                    </div>
                )}
            </Card>
        </motion.div>
    );
}
