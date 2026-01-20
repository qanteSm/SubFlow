/**
 * Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
 * This source code is proprietary. Confidential and private.
 * Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
 */

import axios from 'axios';

// Create axios instance with base configuration
const api = axios.create({
    baseURL: '/api/v1',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Request interceptor for auth token
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('subflow-auth');
    if (token) {
        try {
            const parsed = JSON.parse(token);
            if (parsed.state?.token) {
                config.headers.Authorization = `Bearer ${parsed.state.token}`;
            }
        } catch {
            // Invalid token format
        }
    }
    return config;
});

// Response interceptor for error handling
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Clear auth state on unauthorized
            localStorage.removeItem('subflow-auth');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

// Types
export interface Project {
    id: string;
    tenant_id: string;
    name: string;
    code: string;
    description: string;
    status: 'DRAFT' | 'ACTIVE' | 'ON_HOLD' | 'COMPLETED' | 'CANCELLED';
    contract_amount: number;
    currency: string;
    start_date: string;
    estimated_end_date: string;
    labor_retainage_rate: number;
    material_retainage_rate: number;
    created_at: string;
    updated_at: string;
}

export interface Transaction {
    id: string;
    project_id: string;
    type: 'INVOICE' | 'PAYMENT' | 'RETAINAGE_HELD' | 'RETAINAGE_RELEASE';
    amount_cents: number;
    currency: string;
    effective_date: string;
    description: string;
    reference_no: string;
    created_at: string;
}

export interface AIABillingInput {
    original_contract_sum: number;
    approved_change_orders: number;
    previous_work_completed: number;
    current_work_completed: number;
    stored_materials: number;
    previous_certificates: number;
    labor_retainage_rate: number;
    material_retainage_rate: number;
}

export interface AIABillingResult {
    contract_sum: number;
    total_work_completed: number;
    total_completed_and_stored: number;
    labor_retainage: number;
    material_retainage: number;
    total_retainage: number;
    total_earned: number;
    less_previous_certs: number;
    current_payment_due: number;
    percent_complete: number;
    balance_to_finish: number;
}

export interface FinancialSummary {
    project_id: string;
    total_invoiced: number;
    total_paid: number;
    total_retained: number;
    current_balance: number;
    currency: string;
    transaction_count: number;
}

// API Functions
export const projectsApi = {
    list: () => api.get<{ data: Project[] }>('/projects'),
    get: (id: string) => api.get<Project>(`/projects/${id}`),
    create: (data: Partial<Project>) => api.post<Project>('/projects', data),
    update: (id: string, data: Partial<Project>) =>
        api.put<Project>(`/projects/${id}`, data),
    delete: (id: string) => api.delete(`/projects/${id}`),
    getFinancials: (id: string) =>
        api.get<FinancialSummary>(`/projects/${id}/financials/summary`),
};

export const transactionsApi = {
    listByProject: (projectId: string) =>
        api.get<{ data: Transaction[] }>(`/transactions/project/${projectId}`),
    createInvoice: (data: {
        project_id: string;
        amount: number;
        currency: string;
        invoice_no: string;
    }) => api.post('/transactions/invoice', data),
    createPayment: (data: {
        project_id: string;
        amount: number;
        currency: string;
        bank_receipt_no: string;
    }) => api.post('/transactions/payment', data),
};

export const calculatorApi = {
    calculate: (input: AIABillingInput) =>
        api.post<{ result: AIABillingResult; formatted: Record<string, string> }>(
            '/calculate/aia',
            input
        ),
};

export const systemApi = {
    health: () => api.get('/health'),
    version: () =>
        api.get<{
            application: string;
            version: string;
            architect: string;
            build_id: string;
        }>('/system/version'),
};

export default api;
