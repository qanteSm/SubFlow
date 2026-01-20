/**
 * Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
 * This source code is proprietary. Confidential and private.
 * Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
 */

import { useState } from 'react';
import { Plus, Search, Filter, MoreHorizontal } from 'lucide-react';

// Mock data - in production would come from API
const mockProjects = [
    {
        id: '1',
        name: 'Metro Station Expansion',
        code: 'PRJ-2026-001',
        status: 'ACTIVE',
        contractAmount: 2500000000,
        percentComplete: 45,
        currency: 'TRY',
    },
    {
        id: '2',
        name: 'Hospital Wing Construction',
        code: 'PRJ-2026-002',
        status: 'ACTIVE',
        contractAmount: 1800000000,
        percentComplete: 72,
        currency: 'TRY',
    },
    {
        id: '3',
        name: 'Shopping Mall Renovation',
        code: 'PRJ-2026-003',
        status: 'ON_HOLD',
        contractAmount: 950000000,
        percentComplete: 35,
        currency: 'TRY',
    },
    {
        id: '4',
        name: 'Office Tower Project',
        code: 'PRJ-2025-015',
        status: 'COMPLETED',
        contractAmount: 4200000000,
        percentComplete: 100,
        currency: 'TRY',
    },
];

const statusColors: Record<string, string> = {
    DRAFT: 'bg-gray-100 text-gray-700',
    ACTIVE: 'bg-green-100 text-green-700',
    ON_HOLD: 'bg-yellow-100 text-yellow-700',
    COMPLETED: 'bg-blue-100 text-blue-700',
    CANCELLED: 'bg-red-100 text-red-700',
};

function formatCurrency(cents: number, currency: string): string {
    const major = cents / 100;
    return new Intl.NumberFormat('tr-TR', {
        style: 'currency',
        currency: currency,
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(major);
}

export default function Projects() {
    const [searchQuery, setSearchQuery] = useState('');

    const filteredProjects = mockProjects.filter(
        (project) =>
            project.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            project.code.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Projects</h1>
                    <p className="text-muted-foreground mt-1">
                        Manage your construction projects
                    </p>
                </div>
                <button className="flex items-center gap-2 bg-primary text-primary-foreground px-4 py-2 rounded-lg hover:bg-primary/90 transition-colors">
                    <Plus className="w-5 h-5" />
                    New Project
                </button>
            </div>

            {/* Filters */}
            <div className="flex items-center gap-4">
                <div className="relative flex-1 max-w-md">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
                    <input
                        type="text"
                        placeholder="Search projects..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="w-full pl-10 pr-4 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary"
                    />
                </div>
                <button className="flex items-center gap-2 px-4 py-2 border rounded-lg hover:bg-accent transition-colors">
                    <Filter className="w-5 h-5" />
                    Filters
                </button>
            </div>

            {/* Projects Table */}
            <div className="bg-card rounded-xl border overflow-hidden">
                <table className="w-full">
                    <thead className="bg-muted/50">
                        <tr>
                            <th className="text-left px-6 py-4 text-sm font-medium text-muted-foreground">
                                Project
                            </th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-muted-foreground">
                                Status
                            </th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-muted-foreground">
                                Contract Amount
                            </th>
                            <th className="text-left px-6 py-4 text-sm font-medium text-muted-foreground">
                                Progress
                            </th>
                            <th className="text-right px-6 py-4 text-sm font-medium text-muted-foreground">
                                Actions
                            </th>
                        </tr>
                    </thead>
                    <tbody className="divide-y">
                        {filteredProjects.map((project) => (
                            <tr key={project.id} className="hover:bg-muted/30 transition-colors">
                                <td className="px-6 py-4">
                                    <div>
                                        <p className="font-medium">{project.name}</p>
                                        <p className="text-sm text-muted-foreground">
                                            {project.code}
                                        </p>
                                    </div>
                                </td>
                                <td className="px-6 py-4">
                                    <span
                                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${statusColors[project.status]
                                            }`}
                                    >
                                        {project.status.replace('_', ' ')}
                                    </span>
                                </td>
                                <td className="px-6 py-4 font-medium">
                                    {formatCurrency(project.contractAmount, project.currency)}
                                </td>
                                <td className="px-6 py-4">
                                    <div className="flex items-center gap-3">
                                        <div className="flex-1 h-2 bg-muted rounded-full overflow-hidden max-w-[120px]">
                                            <div
                                                className="h-full bg-primary rounded-full transition-all"
                                                style={{ width: `${project.percentComplete}%` }}
                                            />
                                        </div>
                                        <span className="text-sm text-muted-foreground w-12">
                                            {project.percentComplete}%
                                        </span>
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-right">
                                    <button className="p-2 hover:bg-muted rounded-lg transition-colors">
                                        <MoreHorizontal className="w-5 h-5" />
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>

                {filteredProjects.length === 0 && (
                    <div className="p-12 text-center text-muted-foreground">
                        <p>No projects found</p>
                    </div>
                )}
            </div>

            {/* Pagination */}
            <div className="flex items-center justify-between">
                <p className="text-sm text-muted-foreground">
                    Showing {filteredProjects.length} of {mockProjects.length} projects
                </p>
                <div className="flex items-center gap-2">
                    <button className="px-4 py-2 border rounded-lg hover:bg-accent transition-colors disabled:opacity-50">
                        Previous
                    </button>
                    <button className="px-4 py-2 border rounded-lg hover:bg-accent transition-colors">
                        Next
                    </button>
                </div>
            </div>
        </div>
    );
}
