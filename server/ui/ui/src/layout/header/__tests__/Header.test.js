// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
// component
import Header from '../Header';


test('renders learn react link', () => {
  render(
    <Router>
      <Header />
    </Router>
  );
  const linkElement = screen.getByAltText(/Gigantum HOS Logo/i);
  expect(linkElement).toBeInTheDocument();
});
