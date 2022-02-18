// vendor
import React, { FC } from 'react';
import { Link } from "react-router-dom";
// Images
import logoSVG from 'Images/logos/hos-logo-white-simple.svg';
import User from './user/User';
// css
import './Logo.scss';

interface Props {
  send: any;
  serverName: string;
}

const logoURL = `${window.location.protocol}//${window.location.hostname}/ui/logo.svg`;
let logoText = '';

fetch(logoURL)
    .then(r => r.text())
    .then(text => {
        logoText = text;
    })

const Logo: FC<Props> = ({ send, serverName }: Props) => {

  return (
    <div className="Logo">
      <div className="column-1-span-1 flex justify--right Logo__container">
        <Link
          className="flex"
          to = "/"
        >
          <div className="Logo__img" dangerouslySetInnerHTML={{__html: logoText}} />
        </Link>
      </div>
      <Link
        to="/"
        className="column-1-span-9 Logo__hostname flex align-items--center"
      >
          { serverName }
      </Link>
      <div className="Logo__divider"/>
      <div className="column-1-span-3 Logo__user">
        <User send={send} />
      </div>
    </div>
  )
}


export default Logo;
