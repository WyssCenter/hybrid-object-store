// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import Card from '../Card';

describe('Card', () => {
  it('renders namespace Card', () => {
    render(
      <Card
        description="this is a namespace description"
        name="namespace"
        path="data/namepsace"
      />
    );
    const linkElement = screen.getByText("this is a namespace description");
    expect(linkElement).toBeInTheDocument();
  });


  it('renders dataset Card', () => {
    render(
      <Card
        description="this is a dataset description"
        name="dataset"
        path="data/namepsace/dataset"
      />
    );
    const linkElement = screen.getByText("this is a dataset description");
    expect(linkElement).toBeInTheDocument();
  });
})
