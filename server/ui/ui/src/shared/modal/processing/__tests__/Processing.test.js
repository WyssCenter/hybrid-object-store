import React from 'react';
import { render, screen } from '@testing-library/react';
import Processing from '../Processing';


describe('Processing Create Modal', () => {
  it('renders component', () => {
    render(
      <Processing
        modalType="namespace"
        name="namespace-12"
      />
    );
    const processingElement = screen.getByText(/namespace-12/);
    expect(processingElement).toBeInTheDocument();
  });
})
