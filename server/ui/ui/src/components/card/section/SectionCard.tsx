// vendor
import React, { FC, ReactNode } from 'react';
import classNames from 'classnames';
// css
import '../default/Card.scss';
import './SectionCard.scss';


interface Props {
  children: ReactNode | Element;
  verticalHeight?: string | null;
  noPadding?: boolean | null;
  span?: number | null;
}

const SectionCard: FC<Props> = ({
  children,
  verticalHeight,
  noPadding,
  span,
}: Props) => {

  const sectionCardCSS = classNames({
    'SectionCard Card column-1-span-12': !span,
   [ `SectionCard Card column-1-span-${span}`]: span,
    [verticalHeight]: verticalHeight !== null,
    'SectionCard--noPadding': noPadding,
  });

  return (
    <div className={sectionCardCSS}>
      {children}
    </div>
  )
}


export default SectionCard;
