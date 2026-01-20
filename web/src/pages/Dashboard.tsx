/**
 * Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
 * This source code is proprietary. Confidential and private.
 * Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
 */

import { useQuery } from '@tanstack/react-query';
import {
    TrendingUp,
    TrendingDown,
    DollarSign,
    FolderKanban,
    Receipt,
    PiggyBank,
} from 'lucide-react';
import { systemApi } from '../lib/api';

// Stat Card Component
function StatCard({
    title,
    value,
    change,
    changeType,
    icon: Icon,
}: {
    title: string;
    value: string;
    change: string;
    changeType: 'positive' | 'negative' | 'neutral';
    icon: React.ElementType;
}) {
    return (
        <div className="bg-card rounded-xl border p-6 card-hover">
            <div className="flex items-center justify-between">
                <div>
                    <p className="text-sm text-muted-foreground">{title}</p>
                    <p className="text-2xl font-bold mt-1">{value}</p>
                    <div className="flex items-center gap-1 mt-2">
                        {changeType === 'positive' && (
                            <TrendingUp className="w-4 h-4 text-green-500" />
                        )}
                        {changeType === 'negative' && (
                            <TrendingDown className="w-4 h-4 text-red-500" />
                        )}
                        <span
                            className={`text-sm ${changeType === 'positive'
                                    ? 'text-green-500'
                                    : changeType === 'negative'
                                        ? 'text-red-500'
                                        : 'text-muted-foreground'
                                }`}
                        >
                            {change}
                        </span>
                    </div>
                </div>
                <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center">
                    <Icon className="w-6 h-6 text-primary" />
                </div>
            </div>
        </div>
    );
}

export default function Dashboard() {
    // Fetch system info
    const { data: versionData } = useQuery({
        queryKey: ['system-version'],
        queryFn: () => systemApi.version(),
        staleTime: Infinity,
    });

    // Mock stats - in production would come from API
    const stats = [
        {
            title: 'Total Revenue',
            value: '₺2,450,000',
            change: '+12.5% from last month',
            changeType: 'positive' as const,
            icon: DollarSign,
        },
        {
            title: 'Active Projects',
            value: '24',
            change: '+3 new this month',
            changeType: 'positive' as const,
            icon: FolderKanban,
        },
        {
            title: 'Pending Invoices',
            value: '₺450,000',
            change: '8 awaiting payment',
            changeType: 'neutral' as const,
            icon: Receipt,
        },
        {
            title: 'Retainage Held',
            value: '₺180,000',
            change: '-₺25,000 released',
            changeType: 'negative' as const,
            icon: PiggyBank,
        },
    ];

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
                    <p className="text-muted-foreground mt-1">
                        Enterprise Construction Financial Overview
                    </p>
                </div>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <span>v{versionData?.data?.version || '1.0.0'}</span>
                    <span>•</span>
                    <span>{versionData?.data?.architect || 'SubFlow'}</span>
                </div>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {stats.map((stat) => (
                    <StatCard key={stat.title} {...stat} />
                ))}
            </div>

            {/* Charts Section */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Revenue Chart Placeholder */}
                <div className="bg-card rounded-xl border p-6">
                    <h3 className="font-semibold mb-4">Monthly Revenue</h3>
                    <div className="h-64 flex items-center justify-center text-muted-foreground">
                        <div className="text-center">
                            <div className="text-4xl mb-2">📊</div>
                            <p>Recharts integration ready</p>
                            <p className="text-sm">Connect to real data</p>
                        </div>
                    </div>
                </div>

                {/* Project Status Placeholder */}
                <div className="bg-card rounded-xl border p-6">
                    <h3 className="font-semibold mb-4">Project Status Distribution</h3>
                    <div className="h-64 flex items-center justify-center text-muted-foreground">
                        <div className="text-center">
                            <div className="text-4xl mb-2">🥧</div>
                            <p>Pie chart ready</p>
                            <p className="text-sm">Active / Completed / On Hold</p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Recent Transactions */}
            <div className="bg-card rounded-xl border">
                <div className="p-6 border-b">
                    <h3 className="font-semibold">Recent Transactions</h3>
                </div>
                <div className="divide-y">
                    {[
                        {
                            type: 'INVOICE',
                            project: 'Metro Station Project',
                            amount: '₺125,000',
                            date: '2 hours ago',
                        },
                        {
                            type: 'PAYMENT',
                            project: 'Hospital Expansion',
                            amount: '₺450,000',
                            date: '5 hours ago',
                        },
                        {
                            type: 'RETAINAGE',
                            project: 'Shopping Mall',
                            amount: '₺35,000',
                            date: '1 day ago',
                        },
                    ].map((tx, i) => (
                        <div key={i} className="p-4 flex items-center justify-between">
                            <div className="flex items-center gap-4">
                                <div
                                    className={`w-10 h-10 rounded-full flex items-center justify-center ${tx.type === 'INVOICE'
                                            ? 'bg-blue-100 text-blue-600'
                                            : tx.type === 'PAYMENT'
                                                ? 'bg-green-100 text-green-600'
                                                : 'bg-yellow-100 text-yellow-600'
                                        }`}
                                >
                                    {tx.type === 'INVOICE' ? (
                                        <Receipt className="w-5 h-5" />
                                    ) : tx.type === 'PAYMENT' ? (
                                        <DollarSign className="w-5 h-5" />
                                    ) : (
                                        <PiggyBank className="w-5 h-5" />
                                    )}
                                </div>
                                <div>
                                    <p className="font-medium">{tx.project}</p>
                                    <p className="text-sm text-muted-foreground">{tx.type}</p>
                                </div>
                            </div>
                            <div className="text-right">
                                <p className="font-medium">{tx.amount}</p>
                                <p className="text-sm text-muted-foreground">{tx.date}</p>
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Footer */}
            <div className="text-center text-sm text-muted-foreground pt-8 border-t">
                <p>
                    SubFlow Enterprise Construction Financial Ledger © 2026 Muhammet Ali
                    Büyük
                </p>
                <p className="mt-1">
                    <a
                        href="https://alibuyuk.net"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="hover:text-primary"
                    >
                        alibuyuk.net
                    </a>
                </p>
            </div>
        </div>
    );
}
