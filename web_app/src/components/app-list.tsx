
'use client';

import type { SerializableApp } from '@/lib/types';
import { useState, useMemo } from 'react';
import Link from 'next/link';
import { Input } from '@/components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { getIconComponent } from '@/lib/icon-map';
import { Search } from 'lucide-react';
import { Badge } from '@/components/ui/badge';

type AppListProps = {
  apps: SerializableApp[];
};

export function AppList({ apps }: AppListProps) {
  const [searchTerm, setSearchTerm] = useState('');

  const filteredApps = useMemo(() => {
    if (!searchTerm) return apps;
    return apps.filter(
      app =>
        app.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        app.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
        app.client.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [apps, searchTerm]);

  return (
    <div>
      <div className="relative mb-8">
        <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
        <Input
          type="text"
          placeholder="Search applications..."
          value={searchTerm}
          onChange={e => setSearchTerm(e.target.value)}
          className="w-full max-w-lg pl-12 h-12 text-base rounded-full bg-card shadow-sm"
        />
      </div>

      {filteredApps.length > 0 ? (
        <div className="border rounded-2xl overflow-hidden">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[80px]"></TableHead>
                <TableHead>Application</TableHead>
                <TableHead>Latest Version</TableHead>
                <TableHead>Client</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredApps.map(app => {
                const Icon = getIconComponent(app.iconId);
                return (
                  <TableRow key={app.id} className="cursor-pointer">
                    <TableCell>
                      <Link href={`/apps/${app.id}`} tabIndex={-1}>
                        <div className="p-2 bg-card shadow-sm rounded-lg inline-block">
                          <Icon className="h-8 w-8 text-primary" />
                        </div>
                      </Link>
                    </TableCell>
                    <TableCell className="font-medium">
                      <Link href={`/apps/${app.id}`} className="hover:underline">
                        {app.name}
                      </Link>
                      <p className="text-sm text-muted-foreground font-normal line-clamp-1">{app.description}</p>
                    </TableCell>
                    <TableCell>{app.versions[0].version}</TableCell>
                    <TableCell>
                        <Badge variant="outline">{app.client}</Badge>
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </div>
      ) : (
        <div className="text-center py-16">
          <p className="text-muted-foreground">No applications found matching your search.</p>
        </div>
      )}
    </div>
  );
}
