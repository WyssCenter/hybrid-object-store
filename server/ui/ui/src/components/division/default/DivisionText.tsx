// vendor
import React, { FC } from 'react';
// css
import './DivisionText.scss';


interface Props {
  text: string,
}

const DivisionText: FC<Props> = ({
  text,
}: Props) => {

  return (
    <div
      className="DivisionText"
    >
      <div className="DivisionText__separator" />
      <h5 className="DivisionText__h5">{text}</h5>
      <div className="DivisionText__separator" />
    </div>
  )
}


export default DivisionText;
