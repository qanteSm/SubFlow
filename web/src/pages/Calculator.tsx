/**
 * Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
 * This source code is proprietary. Confidential and private.
 * Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
 * 
 * AIA G702/G703 Billing Calculator Interface
 */

import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Calculator as CalcIcon, RefreshCw, Download } from 'lucide-react';
import { calculatorApi, AIABillingInput } from '../lib/api';

interface CalculatorForm {
    originalContractSum: string;
    approvedChangeOrders: string;
    previousWorkCompleted: string;
    currentWorkCompleted: string;
    storedMaterials: string;
    previousCertificates: string;
    laborRetainageRate: string;
    materialRetainageRate: string;
}

const initialForm: CalculatorForm = {
    originalContractSum: '1000000',
    approvedChangeOrders: '50000',
    previousWorkCompleted: '300000',
    currentWorkCompleted: '150000',
    storedMaterials: '50000',
    previousCertificates: '250000',
    laborRetainageRate: '10',
    materialRetainageRate: '5',
};

function formatCurrency(cents: number): string {
    return new Intl.NumberFormat('tr-TR', {
        style: 'currency',
        currency: 'TRY',
        minimumFractionDigits: 2,
    }).format(cents / 100);
}

export default function Calculator() {
    const [form, setForm] = useState<CalculatorForm>(initialForm);
    const [result, setResult] = useState<{
        result: Record<string, number>;
        formatted: Record<string, string>;
    } | null>(null);

    const mutation = useMutation({
        mutationFn: (input: AIABillingInput) => calculatorApi.calculate(input),
        onSuccess: (response) => {
            setResult(response.data);
        },
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        // Convert form values to cents and basis points
        const input: AIABillingInput = {
            original_contract_sum: Math.round(parseFloat(form.originalContractSum) * 100),
            approved_change_orders: Math.round(parseFloat(form.approvedChangeOrders) * 100),
            previous_work_completed: Math.round(parseFloat(form.previousWorkCompleted) * 100),
            current_work_completed: Math.round(parseFloat(form.currentWorkCompleted) * 100),
            stored_materials: Math.round(parseFloat(form.storedMaterials) * 100),
            previous_certificates: Math.round(parseFloat(form.previousCertificates) * 100),
            labor_retainage_rate: Math.round(parseFloat(form.laborRetainageRate) * 100), // 10% -> 1000 basis points
            material_retainage_rate: Math.round(parseFloat(form.materialRetainageRate) * 100),
        };

        mutation.mutate(input);
    };

    const handleReset = () => {
        setForm(initialForm);
        setResult(null);
    };

    const handleChange = (field: keyof CalculatorForm, value: string) => {
        setForm((prev) => ({ ...prev, [field]: value }));
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">AIA Calculator</h1>
                    <p className="text-muted-foreground mt-1">
                        G702/G703 Application & Certificate for Payment
                    </p>
                </div>
                <div className="flex items-center gap-2">
                    <button
                        onClick={handleReset}
                        className="flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-accent transition-colors"
                    >
                        <RefreshCw className="w-5 h-5" />
                        Reset
                    </button>
                    <button className="flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-accent transition-colors">
                        <Download className="w-5 h-5" />
                        Export PDF
                    </button>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Input Form */}
                <div className="bg-card rounded-xl border p-6">
                    <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
                        <CalcIcon className="w-5 h-5" />
                        Input Values
                    </h2>

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">
                                    Original Contract Sum (₺)
                                </label>
                                <input
                                    type="number"
                                    value={form.originalContractSum}
                                    onChange={(e) =>
                                        handleChange('originalContractSum', e.target.value)
                                    }
                                    className="w-full px-3 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium mb-1">
                                    Approved Change Orders (₺)
                                </label>
                                <input
                                    type="number"
                                    value={form.approvedChangeOrders}
                                    onChange={(e) =>
                                        handleChange('approvedChangeOrders', e.target.value)
                                    }
                                    className="w-full px-3 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                                />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">
                                    Previous Work Completed (₺)
                                </label>
                                <input
                                    type="number"
                                    value={form.previousWorkCompleted}
                                    onChange={(e) =>
                                        handleChange('previousWorkCompleted', e.target.value)
                                    }
                                    className="w-full px-3 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium mb-1">
                                    Current Work Completed (₺)
                                </label>
                                <input
                                    type="number"
                                    value={form.currentWorkCompleted}
                                    onChange={(e) =>
                                        handleChange('currentWorkCompleted', e.target.value)
                                    }
                                    className="w-full px-3 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                                />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">
                                    Stored Materials (₺)
                                </label>
                                <input
                                    type="number"
                                    value={form.storedMaterials}
                                    onChange={(e) =>
                                        handleChange('storedMaterials', e.target.value)
                                    }
                                    className="w-full px-3 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium mb-1">
                                    Previous Certificates (₺)
                                </label>
                                <input
                                    type="number"
                                    value={form.previousCertificates}
                                    onChange={(e) =>
                                        handleChange('previousCertificates', e.target.value)
                                    }
                                    className="w-full px-3 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                                />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">
                                    Labor Retainage Rate (%)
                                </label>
                                <input
                                    type="number"
                                    step="0.1"
                                    value={form.laborRetainageRate}
                                    onChange={(e) =>
                                        handleChange('laborRetainageRate', e.target.value)
                                    }
                                    className="w-full px-3 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium mb-1">
                                    Material Retainage Rate (%)
                                </label>
                                <input
                                    type="number"
                                    step="0.1"
                                    value={form.materialRetainageRate}
                                    onChange={(e) =>
                                        handleChange('materialRetainageRate', e.target.value)
                                    }
                                    className="w-full px-3 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                                />
                            </div>
                        </div>

                        <button
                            type="submit"
                            disabled={mutation.isPending}
                            className="w-full bg-primary text-primary-foreground py-3 rounded-lg hover:bg-primary/90 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
                        >
                            {mutation.isPending ? (
                                <>
                                    <RefreshCw className="w-5 h-5 animate-spin" />
                                    Calculating...
                                </>
                            ) : (
                                <>
                                    <CalcIcon className="w-5 h-5" />
                                    Calculate
                                </>
                            )}
                        </button>
                    </form>
                </div>

                {/* Results */}
                <div className="bg-card rounded-xl border p-6">
                    <h2 className="text-lg font-semibold mb-4">Calculation Results</h2>

                    {result ? (
                        <div className="space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <div className="bg-muted/50 rounded-lg p-4">
                                    <p className="text-sm text-muted-foreground">Contract Sum</p>
                                    <p className="text-xl font-bold">
                                        {formatCurrency(result.result.contract_sum)}
                                    </p>
                                </div>
                                <div className="bg-muted/50 rounded-lg p-4">
                                    <p className="text-sm text-muted-foreground">% Complete</p>
                                    <p className="text-xl font-bold">
                                        {(result.result.percent_complete / 100).toFixed(1)}%
                                    </p>
                                </div>
                            </div>

                            <div className="border-t pt-4 space-y-3">
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">
                                        Total Work Completed
                                    </span>
                                    <span className="font-medium">
                                        {formatCurrency(result.result.total_work_completed)}
                                    </span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">
                                        + Stored Materials
                                    </span>
                                    <span className="font-medium">
                                        {formatCurrency(
                                            result.result.total_completed_and_stored -
                                            result.result.total_work_completed
                                        )}
                                    </span>
                                </div>
                                <div className="flex justify-between border-t pt-2">
                                    <span className="text-muted-foreground">
                                        Total Completed & Stored
                                    </span>
                                    <span className="font-semibold">
                                        {formatCurrency(result.result.total_completed_and_stored)}
                                    </span>
                                </div>
                            </div>

                            <div className="border-t pt-4 space-y-3">
                                <div className="flex justify-between text-yellow-600">
                                    <span>Labor Retainage (10%)</span>
                                    <span>
                                        -{formatCurrency(result.result.labor_retainage)}
                                    </span>
                                </div>
                                <div className="flex justify-between text-yellow-600">
                                    <span>Material Retainage (5%)</span>
                                    <span>
                                        -{formatCurrency(result.result.material_retainage)}
                                    </span>
                                </div>
                                <div className="flex justify-between border-t pt-2">
                                    <span className="text-muted-foreground">Total Retainage</span>
                                    <span className="font-semibold text-yellow-600">
                                        -{formatCurrency(result.result.total_retainage)}
                                    </span>
                                </div>
                            </div>

                            <div className="border-t pt-4 space-y-3">
                                <div className="flex justify-between">
                                    <span className="text-muted-foreground">Total Earned</span>
                                    <span className="font-medium">
                                        {formatCurrency(result.result.total_earned)}
                                    </span>
                                </div>
                                <div className="flex justify-between text-red-600">
                                    <span>Less Previous Certificates</span>
                                    <span>
                                        -{formatCurrency(result.result.less_previous_certs)}
                                    </span>
                                </div>
                            </div>

                            <div className="bg-primary/10 rounded-lg p-4 border-2 border-primary">
                                <div className="flex justify-between items-center">
                                    <span className="font-semibold text-lg">
                                        Current Payment Due
                                    </span>
                                    <span className="text-2xl font-bold text-primary">
                                        {formatCurrency(result.result.current_payment_due)}
                                    </span>
                                </div>
                            </div>

                            <div className="flex justify-between text-muted-foreground text-sm pt-2">
                                <span>Balance to Finish</span>
                                <span>{formatCurrency(result.result.balance_to_finish)}</span>
                            </div>
                        </div>
                    ) : (
                        <div className="h-64 flex items-center justify-center text-muted-foreground">
                            <div className="text-center">
                                <CalcIcon className="w-12 h-12 mx-auto mb-4 opacity-50" />
                                <p>Enter values and click Calculate</p>
                                <p className="text-sm mt-1">Results will appear here</p>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Info */}
            <div className="bg-blue-50 dark:bg-blue-950/20 rounded-xl p-6 border border-blue-200 dark:border-blue-800">
                <h3 className="font-semibold text-blue-800 dark:text-blue-200 mb-2">
                    About AIA G702/G703
                </h3>
                <p className="text-sm text-blue-700 dark:text-blue-300">
                    The AIA G702 (Application and Certificate for Payment) and G703
                    (Continuation Sheet) are standard forms used in the construction
                    industry for contractors to request progress payments. This calculator
                    uses BigInt arithmetic to ensure precise financial calculations
                    without floating-point errors.
                </p>
                <p className="text-xs text-blue-600 dark:text-blue-400 mt-2">
                    Architect: Muhammet Ali Büyük • alibuyuk.net
                </p>
            </div>
        </div>
    );
}
