import React, { Suspense } from "react"
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
// import HomePage from './pages/Home';
// import StatsPage from './pages/Stats';
// import AboutPage from './pages/About';

const AboutPage = React.lazy(() => import('./pages/About'))
const StatsPage = React.lazy(() => import('./pages/Stats'))
const HomePage = React.lazy(() => import('./pages/Home'))

function App() {
  return (
    <Router>
      <Suspense fallback={<div>Loading...</div>}>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/stats" element={<StatsPage />} />
        <Route path="/about" element={<AboutPage/>} />
      </Routes>
      </Suspense>
    </Router>
  );
}

export default App
