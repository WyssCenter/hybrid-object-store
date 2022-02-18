// vendor
import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import Search from '../Search';

describe('Search', () => {

  it('Renders Search', () => {
    render(
      <Search
        list={[]}
        placeholder="Search Datasets"
        updateList={jest.fn()}
      />
    );
    const labelElement = screen.getByText('Search');
    const textboxElement = screen.getByRole('textbox');
    expect(labelElement).toBeInTheDocument();
    expect(textboxElement.placeholder).toBe('Search Datasets');
  });



  it('calls update function', () => {
    const updateList = jest.fn();
    render(
      <Search
        list={[]}
        placeholder="Search Datasets"
        updateList={updateList}
      />
    );
    // first call on mount
    expect(updateList).toHaveBeenCalledTimes(1);
    const textboxElement = screen.getByRole('textbox');
    fireEvent.keyUp(textboxElement, { target: { value: 'sample input' }})

    expect(textboxElement.value).toBe('sample input');
    expect(updateList).toHaveBeenCalledTimes(2);
  });


  it('calls updateList with original list by default', () => {
    const updateList = jest.fn();
    const sampleList = [
      { name: 'test1' },
      { name: 'test2' },
      { name: 'test3' },
    ]
    render(
      <Search
        list={sampleList}
        placeholder="Search Datasets"
        updateList={updateList}
      />
    );
    expect(updateList.mock.calls[0][0].length).toEqual(3);
  });

  it('calls updateList with the correct number of matches', () => {
    const updateList = jest.fn();
    const sampleList = [
      { name: 'test1' },
      { name: 'test1' },
      { name: 'test3' },
    ]
    render(
      <Search
        list={sampleList}
        placeholder="Search Datasets"
        updateList={updateList}
      />
    );
    const textboxElement = screen.getByRole('textbox');
    fireEvent.keyUp(textboxElement, { target: { value: 'test1' }})
    expect(textboxElement.value).toBe('test1');
    expect(updateList.mock.calls[1][0].length).toEqual(2);
  });

  it('calls updateList with empty array when no match', () => {
    const updateList = jest.fn();
    const sampleList = [
      { name: 'test1' },
      { name: 'test2' },
      { name: 'test3' },
    ]
    render(
      <Search
        list={sampleList}
        placeholder="Search Datasets"
        updateList={updateList}
      />
    );
    const textboxElement = screen.getByRole('textbox');
    fireEvent.keyUp(textboxElement, { target: { value: '---' }})
    expect(textboxElement.value).toBe('---');
    expect(updateList.mock.calls[1][0].length).toEqual(0);
  });
})
