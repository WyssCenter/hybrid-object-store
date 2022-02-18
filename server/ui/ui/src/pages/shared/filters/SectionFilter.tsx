// vendor
import React, { FC } from 'react';
// components
import { Search } from 'Components/form/text/index';
// style
import './SectionFilter.scss';


interface Props {
  list: any;
  formattedSection: string;
  updateList: (list: any) => void,
}


const SectionFilter: FC<Props> = ({
  list,
  formattedSection,
  updateList,
}:Props) => {
  return (
    <section className="SectionFilter">
      <Search
        disabled={false}
        placeholder={`Search ${formattedSection}s`}
        list={list}
        updateList={updateList}
      />
    </section>
  )
}

export default SectionFilter;
