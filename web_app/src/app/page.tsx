import { getApps } from '@/lib/data';
import type { SerializableApp } from '@/lib/types';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { AppList } from '@/components/app-list';

export default async function Home() {
  const apps = await getApps();

  const serializableApps: SerializableApp[] = apps.map(({ icon, ...rest }) => rest);

  return (
    <div className="container mx-auto px-4 py-8">
      <Tabs defaultValue="all" className="mt-8">
        <TabsList className="mb-6">
          <TabsTrigger value="all">All Apps</TabsTrigger>
          <TabsTrigger value="productivity">Productivity</TabsTrigger>
          <TabsTrigger value="finance">Finance</TabsTrigger>
          <TabsTrigger value="development">Development</TabsTrigger>
        </TabsList>
        <TabsContent value="all">
           <AppList apps={serializableApps} />
        </TabsContent>
        <TabsContent value="productivity">
           <AppList apps={serializableApps.filter(a => a.category === 'Productivity')} />
        </TabsContent>
        <TabsContent value="finance">
           <AppList apps={serializableApps.filter(a => a.category === 'Finance')} />
        </TabsContent>
        <TabsContent value="development">
           <AppList apps={serializableApps.filter(a => a.category === 'Development')} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
