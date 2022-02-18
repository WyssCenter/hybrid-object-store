// vendor
import React from 'react';
import { act } from 'react-dom/test-utils';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// context
import AppContext from 'Src/AppContext';
// components
import SectionFilter from '../SectionFilter';
// import mockNamepacesData from '../../';



describe('SectionFilter', () => {
  beforeAll(() => {
    let modalElement = window.document.createElement('div');
    modalElement.setAttribute('id', 'modal');

    document.body.appendChild(modalElement)
  })

  it('Renders SectionFilter', async () => {


    render(
      <SectionFilter
        namespaces={[]}
        permissions
        formattedSection="Namespace"
        section="namespace"
      />
    );
    const createButton = screen.getByText(/Create Namespace/i);

    expect(createButton).toBeInTheDocument();

  });

});
