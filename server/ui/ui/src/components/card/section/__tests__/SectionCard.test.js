// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import SectionCard from '../SectionCard';

describe('SectionCard', () => {
  it('Renders "child content"', () => {
    render(
      <SectionCard
        verticalHeight="grid-v-3"
      >
        <p>child content</p>
      </SectionCard>
    );
    const cardElement = screen.getByText("child content");
    expect(cardElement).toBeInTheDocument();
  });


  it('Renders "other child content"', () => {
    render(
      <SectionCard
        verticalHeight="grid-v-6"
      >
        <p>other child content</p>
      </SectionCard>
    );
    const cardElement = screen.getByText("other child content");
    expect(cardElement).toBeInTheDocument();
  });
})
