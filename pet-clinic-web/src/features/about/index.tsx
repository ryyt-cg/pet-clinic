import {Header} from "@/components/layout/header.tsx";
import {Search} from "@/components/search.tsx";
import {ThemeSwitch} from "@/components/theme-switch.tsx";
import {ConfigDrawer} from "@/components/config-drawer.tsx";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const About = () => {
  return (
      <>
        <Header fixed>
          <div className='ms-auto flex items-center space-x-4'>
            <Search />
            <ThemeSwitch />
            <ConfigDrawer />
          </div>
        </Header>

        <div className="container mx-auto py-12">
          {/* Hero Section */}
          <section className="text-center mb-16">
            <h1 className="text-5xl font-bold mb-4">About Pet Clinic</h1>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Rebuild Pet Clinic using
              <ul>
                <li>1) Frontend: React, Vite, TailwindCSS </li>
                <li>2) Backend: Golang, Fiber, GORM</li>
                <li>3) Database: SQLite, PostgreSQL</li>
              </ul>
            </p>
            <Button className="mt-8">Learn More</Button>
          </section>

          {/* My Mission/Values Section */}
          <section className="grid md:grid-cols-2 gap-8 mb-16">
            <Card>
              <CardHeader>
                <CardTitle>My Goal</CardTitle>
              </CardHeader>
              <CardContent>
                <p>My goal is to apply best practices and production code.  The project could be used as a template</p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardTitle>My Values</CardTitle>
              </CardHeader>
              <CardContent>
                <ul>
                  <li>Innovation</li>
                  <li>Integrity</li>
                </ul>
              </CardContent>
            </Card>
          </section>
        </div>
      </>
  );
};

export default About;