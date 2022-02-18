// vendor
import React, { FC } from 'react';
import PatSection from './section/PatSection';
// css
import './PersonalAccessTokens.scss';



const PersonalAccessTokens: FC = () => {

  return (
    <div className="grid">
      <div className="PersonalAccessTokens">
        <PatSection />
      </div>
    </div>
  );
}

export default PersonalAccessTokens;
