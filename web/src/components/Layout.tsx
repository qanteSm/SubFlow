/**
 * Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
 * This source code is proprietary. Confidential and private.
 * Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
 */

import { Outlet, Link, useLocation } from 'react-router-dom';
import {
    LayoutDashboard,
    FolderKanban,
    Calculator,
    Settings,
    Menu,
    Moon,
    Sun,
} from 'lucide-react';
import { useState } from 'react';
import { useThemeStore } from '../store/themeStore';

const navigation = [
    { name: 'Dashboard', href: '/', icon: LayoutDashboard },
    { name: 'Projects', href: '/projects', icon: FolderKanban },
    { name: 'Calculator', href: '/calculator', icon: Calculator },
    { name: 'Settings', href: '/settings', icon: Settings },
];

export default function Layout() {
    const location = useLocation();
    const [sidebarOpen, setSidebarOpen] = useState(true);
    const { isDark, toggle } = useThemeStore();

    return (
        <div className={`min-h-screen ${isDark ? 'dark' : ''}`}>
            <div className="flex h-screen bg-background">
                {/* Sidebar */}
                <aside
                    className={`${sidebarOpen ? 'w-64' : 'w-20'
                        } bg-card border-r transition-all duration-300 flex flex-col`}
                >
                    {/* Logo */}
                    <div className="h-16 flex items-center px-4 border-b">
                        <Link to="/" className="flex items-center gap-2">
                            <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-purple-600 rounded-xl flex items-center justify-center">
                                <span className="text-white font-bold text-lg">SF</span>
                            </div>
                            {sidebarOpen && (
                                <span className="font-semibold text-lg gradient-text">
                                    SubFlow
                                </span>
                            )}
                        </Link>
                    </div>

                    {/* Navigation */}
                    <nav className="flex-1 px-3 py-4 space-y-1">
                        {navigation.map((item) => {
                            const isActive = location.pathname === item.href;
                            return (
                                <Link
                                    key={item.name}
                                    to={item.href}
                                    className={`flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all ${isActive
                                            ? 'bg-primary text-primary-foreground'
                                            : 'text-muted-foreground hover:bg-accent hover:text-foreground'
                                        }`}
                                >
                                    <item.icon className="w-5 h-5" />
                                    {sidebarOpen && <span>{item.name}</span>}
                                </Link>
                            );
                        })}
                    </nav>

                    {/* Footer */}
                    <div className="p-4 border-t">
                        {sidebarOpen && (
                            <div className="text-xs text-muted-foreground text-center">
                                <p>© 2026 Muhammet Ali Büyük</p>
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
                        )}
                    </div>
                </aside>

                {/* Main Content */}
                <div className="flex-1 flex flex-col overflow-hidden">
                    {/* Header */}
                    <header className="h-16 bg-card border-b flex items-center justify-between px-6">
                        <button
                            onClick={() => setSidebarOpen(!sidebarOpen)}
                            className="p-2 rounded-lg hover:bg-accent"
                        >
                            <Menu className="w-5 h-5" />
                        </button>

                        <div className="flex items-center gap-4">
                            <button
                                onClick={toggle}
                                className="p-2 rounded-lg hover:bg-accent"
                            >
                                {isDark ? (
                                    <Sun className="w-5 h-5" />
                                ) : (
                                    <Moon className="w-5 h-5" />
                                )}
                            </button>

                            <div className="flex items-center gap-2">
                                <div className="w-8 h-8 bg-primary rounded-full flex items-center justify-center text-primary-foreground text-sm font-medium">
                                    MA
                                </div>
                            </div>
                        </div>
                    </header>

                    {/* Page Content */}
                    <main className="flex-1 overflow-auto p-6 bg-background">
                        <Outlet />
                    </main>
                </div>
            </div>
        </div>
    );
}
