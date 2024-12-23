import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

test('renders Baggsy header', () => {
  render(<App />);
  const headerElement = screen.getByText(/Baggsy/i);
  expect(headerElement).toBeInTheDocument();
});

test('handles empty bags gracefully', () => {
  render(<App />);
  const listElement = screen.queryByRole('list');
  expect(listElement).toBeEmptyDOMElement();
});
