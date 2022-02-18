// vendor
import React, { FC, ReactNode } from 'react';
// css
import './Card.scss';


interface Props {
  description: string;
  key: string;
  name: string;
  path: string;

}

const Card: FC<Props> = ({
  description,
  key,
  name,
  path,
}: Props) => {

  return (
    <div className="Card column-4-span-3 grid-v-3">
      <h4 className="Card__title">{name}</h4>
      <p className="Card__p--small">{path}</p>
      <p>{description}</p>
    </div>
  )
}


export default Card;
