// vendor
import React, { FC, useState } from 'react';
import classNames from 'classnames';
// css
import { ReactComponent as AscendingIcon } from 'Images/icons/icon-sort-asc.svg';
import { ReactComponent as DescendingIcon } from 'Images/icons/icon-sort-desc.svg';
import './Header.scss';


interface Props {
  setAllSelected: any;
  allSelected: boolean;
  someSelected: boolean;
  selectedSort: { type: string; ascending: boolean; };
  setSelectedSort: (data: any) => void;
  searchMode: boolean;
}


const Header:FC<Props> = ({
  setAllSelected,
  allSelected,
  someSelected,
  selectedSort,
  setSelectedSort,
  searchMode,
}: Props) => {

  const checkboxCSS = classNames({
    FileHeader__checkbox: true,
    'FileHeader__checkbox--selected': allSelected,
    'FileHeader__checkbox--partial': someSelected && !allSelected,
  })

  const fileCSS = classNames({
    FileHeader__name: true,
    'FileHeader__name--file': true,
    'FileHeader--ascending': selectedSort
      && selectedSort.type === 'file'
      && selectedSort.ascending,
    'FileHeader--descending': selectedSort
      && selectedSort.type === 'file'
      && !selectedSort.ascending,
  })

  const sizeCSS = classNames({
    FileHeader__name: true,
    'FileHeader__name--size': true,
    'FileHeader--ascending': selectedSort
      && selectedSort.type === 'size'
      && selectedSort.ascending,
    'FileHeader--descending': selectedSort
      && selectedSort.type === 'size'
      && !selectedSort.ascending,
  })
  const modifiedCSS = classNames({
    FileHeader__name: true,
    'FileHeader__name--modified': true,
    'FileHeader--ascending': selectedSort
      && selectedSort.type === 'modified'
      && selectedSort.ascending,
    'FileHeader--descending': selectedSort
      && selectedSort.type === 'modified'
      && !selectedSort.ascending,
  })

  return (
    <div className="FileHeader">
      {
        searchMode && (
          <div className="FileHeader__checkbox--placeholder" />
        )
      }
      {
        !searchMode && (
        <div
          className={checkboxCSS}
          onClick={setAllSelected}
        />
        )
      }
      <div
        className={fileCSS}
        onClick={() => {
          setSelectedSort((prevState: any) => ({
            type: 'file',
            ascending: prevState.type === 'file'
              ? !prevState.ascending
              : false,
          }))
        }}
      >
        File
      {
        selectedSort.type === 'file' && selectedSort.ascending && (
          <AscendingIcon />
        )
      }
      {
        selectedSort.type === 'file' && !selectedSort.ascending && (
          <DescendingIcon />
        )
      }
      </div>
      <div
        className={sizeCSS}
        onClick={() => {
          setSelectedSort((prevState: any) => ({
            type: 'size',
            ascending: prevState.type === 'size'
              ? !prevState.ascending
              : false,
          }))
        }}
      >
        Size
      {
        selectedSort.type === 'size' && selectedSort.ascending && (
          <AscendingIcon />
        )
      }
      {
        selectedSort.type === 'size' && !selectedSort.ascending && (
          <DescendingIcon />
        )
      }
      </div>
      <div
        className={modifiedCSS}
        onClick={() => {
          setSelectedSort((prevState: any) => ({
            type: 'modified',
            ascending: prevState.type === 'modified'
              ? !prevState.ascending
              : false,
          }))
        }}
      >
        Modified
      {
        selectedSort.type === 'modified' && selectedSort.ascending && (
          <AscendingIcon />
        )
      }
      {
        selectedSort.type === 'modified' && !selectedSort.ascending && (
          <DescendingIcon />
        )
      }
      </div>
      <div className="FileHeader__name FileHeader__name--actions">
        Actions
      </div>
    </div>
  );
}

export default Header;
