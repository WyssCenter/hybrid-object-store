// vendor
import React, { FC } from 'react';
// css
import './HierarchyHeader.scss';

interface Props {
  header: string;
  subheader: string;
}

const HierarchyHeader: FC<Props> = ({ header, subheader }: Props) => {
  return (
    <div className="grid">
      <div className="HierarchyHeader">
        <div className="HierarchyHeader__spacer column-1-span-12 flex align--center">
          <h2 className="HierarchyHeader__h2"><b>{header}</b>: {subheader}</h2>


        </div>
      </div>
    </div>
  );
}


export default HierarchyHeader;
