export interface User {
    id: string;
    email: string;
}

export interface Monitor {
    id: string;
    user_id: string;
    url: string;
    interval: number;
    last_checked: string | null;
    status_code: number;
    response_time: number;
    is_healthy: boolean;
    is_running: boolean;
    ai_explanation?: string;
}

export interface ApiResponse<T> {
    success: boolean;
    message: string;
    data: T;
}
