// vendor
import React, { FC, MouseEvent } from 'react';
import classNames from 'classnames';
import { IconDefinition }from '@fortawesome/fontawesome-common-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
// css
import './FlatIconTextButton.scss';


type ColorText = 'white' | 'primary' | 'primaryText' | 'secondary' | 'primaryText';

type ColorBackground = 'white' | 'primary' | 'primaryText' | 'secondary' | 'primaryText' | 'transparent';


interface Props {
  click: (event: MouseEvent) => void;
  color: ColorText;
  backgroundColor?: ColorBackground;
  disabled?: boolean;
  text: string;
  icon: IconDefinition;
}

interface ColorsData {
  white: string;
  primary: string;
  secondary: string;
  slateBlue: string;
  primaryText: string;
  slateBlueOpaque: string;
  primaryLight: string;
}

const IconButton: FC<Props> = ({
  click,
  color = 'primary',
  disabled = false,
  text,
  icon,
}: Props) => {
  const root = document.documentElement;
  const primaryHash = root.style.cssText ?
    root.style.cssText.split(';')[0].split(':')[1]
    : '#957299';
  const secondaryHash = root.style.cssText ?
    root.style.cssText.split(';')[0].split(':')[1]
    : '#2f8da3';

  const colorsData: ColorsData = {
  'white': '#fefefe',
  'primary': primaryHash,
  'secondary': secondaryHash,
  'slateBlue': '#364454',
  'primaryText': '#0b1425',
  'slateBlueOpaque': 'rgba(54,68,84, 0.6)',
  'primaryLight': 'rgba(39,70,134, 0.5)',
}

  // vars
  const colorSelected = colorsData[color];
  // css
  const iconButtonCSS = classNames({
    "FlatIconTextButton": true,
    [`FlatIconTextButton--${color}`]: color !== null
  });

  return (
    <button
      className={iconButtonCSS}
      disabled={disabled}
      onClick={click}
    >
      <FontAwesomeIcon icon={icon} color={colorSelected} />
      {' '}
      {text}
    </button>
  )
}


export default IconButton;
