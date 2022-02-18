import React from 'react';
import { render, screen } from '@testing-library/react';
// components
import Progress from '../Progress';

describe('ProgressLoader', () => {
  it('renders learn react link', () => {
    render(
      <Progress
      isCanceling={false}
      isComplete={false}
      error=""
      percentageComplete={50}
      text="Updating Permissions"
      />
    );
    const ProgressText = screen.getByText('Updating Permissions');
    expect(ProgressText).toBeInTheDocument();
  });
});
