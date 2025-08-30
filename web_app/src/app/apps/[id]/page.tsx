import { getAppById, getApps } from '@/lib/data';
import { notFound } from 'next/navigation';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { ArrowLeft, Download, Star } from 'lucide-react';
import { format, parseISO } from 'date-fns';
import { getIconComponent } from '@/lib/icon-map';
import { Separator } from '@/components/ui/separator';

type AppDetailPageProps = {
  params: { id: string };
};

export async function generateStaticParams() {
    const apps = await getApps();
    return apps.map(app => ({ id: app.id }));
}

const handleDownload = async (appName: string, version: string) => {
    'use server';
    console.log(`Tracking download for ${appName} version ${version} by user...`);
};

export default async function AppDetailPage({ params }: AppDetailPageProps) {
  const app = await getAppById(params.id);

  if (!app) {
    notFound();
  }
  const Icon = getIconComponent(app.iconId);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-12">
        <Link href="/" className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground mb-6">
          <ArrowLeft className="h-4 w-4" />
          Back to all apps
        </Link>
        
        <div className="flex flex-col md:flex-row items-start gap-8">
          <div className="p-4 bg-card shadow-sm rounded-3xl">
            <Icon className="h-24 w-24 text-primary" />
          </div>
          <div className="flex-1">
              <h1 className="text-4xl font-bold mt-2">{app.name}</h1>
              <p className="text-lg text-primary font-medium mt-1">{app.category}</p>
              <p className="text-base text-muted-foreground mt-4 max-w-prose">
                  {app.description}
              </p>
              <div className="mt-6">
                 <Button size="lg" className="rounded-full text-lg h-12 px-8">
                    <Download className="mr-3 h-5 w-5" />
                    Download Latest
                 </Button>
              </div>
          </div>
        </div>
      </div>
      
      <Separator />

      <div className="my-12">
        <h2 className="text-2xl font-bold mb-4">What's New</h2>
        <div className="p-6 bg-card rounded-2xl">
            <h3 className="font-semibold text-lg">Version {app.versions[0].version}</h3>
            <p className="text-sm text-muted-foreground mb-3">{format(parseISO(app.versions[0].releaseDate), 'MMMM d, yyyy')}</p>
            <p>{app.versions[0].changelog}</p>
        </div>
      </div>
      
      <Separator />

      <div className="my-12">
        <h2 className="text-2xl font-bold mb-4">Version History</h2>
        <div className="border rounded-2xl overflow-hidden">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[120px]">Version</TableHead>
                  <TableHead className="w-[150px]">Release Date</TableHead>
                  <TableHead>Changelog</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {app.versions.map((version) => (
                  <TableRow key={version.version}>
                    <TableCell className="font-medium">{version.version}</TableCell>
                    <TableCell>{format(parseISO(version.releaseDate), 'MMM d, yyyy')}</TableCell>
                    <TableCell>{version.changelog}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
        </div>
      </div>
    </div>
  );
}
