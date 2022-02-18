// vendor
import React from 'react';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
// components
import CardList from '../CardList';

const namespacesList = [{
    bucketName: 'data',
    description: 'this is a namespace',
    name: 'namespace1'
  },
  {
    bucketName: 'data',
    description: 'this is aanother namespace',
    name: 'namespace2'
  }]

describe('CardList', () => {

  it('Renders CardList item 1', () => {
    render(
      <MemoryRouter>
        <CardList
          namespacesList={namespacesList}
          updateNamespaceFetchId={jest.fn()}
        />
      </MemoryRouter>
    );
    const cardTitle1 = screen.getByText('namespace1');
    expect(cardTitle1).toBeInTheDocument();
  });

  it('Renders CardList item 2', () => {
    render(
      <MemoryRouter>
        <CardList
          namespacesList={namespacesList}
          updateNamespaceFetchId={jest.fn()}
        />
      </MemoryRouter>
    );
    const cardTitle2 = screen.getByText('namespace2');
    expect(cardTitle2).toBeInTheDocument();
  });
})
