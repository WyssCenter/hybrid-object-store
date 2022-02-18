// vendor
import React, { FC, MouseEvent } from 'react';
import classNames from 'classnames';
// css
import './ExpandButton.scss';


interface Props {
  click: (event: MouseEvent) => void;
  disabled?: boolean;
  isExpanded: boolean;
  text: string,
}

const ExpandButton: FC<Props> = ({
  click,
  disabled = false,
  isExpanded,
  text,
}: Props) => {

  const expandButtonCSS = classNames({
    ExpandButton: true,
    "ExpandButton--collapsed": !isExpanded,
    "ExpandButton--expanded": isExpanded,
  })


  const expandIconCSS = classNames({
    ExpandButton__icon: true,
    "ExpandButton__icon--collapsed": !isExpanded,
    "ExpandButton__icon--expanded": isExpanded,
  })

  return (
    <div
      className={expandButtonCSS}
      onClick={click}
      role="presentation"
    >
      <div className="ExpandButton__separator" />
      <button
        className="ExpandButton__button"
        disabled={disabled}
      >
        {text}
      </button>

      <div className={expandIconCSS} />

      <div className="ExpandButton__separator" />
    </div>
  )
}


export default ExpandButton;
