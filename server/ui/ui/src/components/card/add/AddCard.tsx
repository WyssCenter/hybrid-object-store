// vendor
import React, { FC } from 'react';
import classNames from 'classnames';
// components
import { SecondaryButton } from 'Components/button/index';
// css
import '../default/Card.scss';
import './AddCard.scss';


interface Props {
  type: string;
  updateModalVisible: (arg: boolean) => void;
}

const AddCard: FC<Props> = ({
  type,
  updateModalVisible,
}: Props) => {

  const cardCSS = classNames({
    "Card AddCard column-4-span-3": true,
    [`Card--${type}`]: true
  });

  return (
    <div
      className={cardCSS}
    >

      <h4 className="Card__title">Add {type}</h4>

      <SecondaryButton
        click={() => updateModalVisible(true)}
        disabled={false}
        text={`Create ${type}`}
      />
    </div>
  )
}


export default AddCard;
