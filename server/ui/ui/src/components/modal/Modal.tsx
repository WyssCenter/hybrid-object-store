// vendor
import React, { FC, ReactNode } from 'react';
import ReactDom from 'react-dom';
import classNames from 'classnames';
// css
import './Modal.scss';

interface Props {
  children: Element | ReactNode;
  handleClose: (event: React.MouseEvent<Element, MouseEvent>) => void;
  header: string;
  icon?: string;
  overflow?: string,
  size?: string;
  subheader?: string;
  noCancel?: boolean;
}

const Modal: FC<Props> = ({
  children,
  handleClose,
  header,
  overflow = 'visible',
  icon = '',
  size = 'medium',
  subheader = '',
  noCancel = false,
}: Props) => {

  const overflowStyle = (overflow === 'visible') ? 'visible' : 'hidden';
  const modalContentCSS = classNames({
    Modal__content: true,
    [`Modal__content--${size}`]: size, // large, medium, small
    [icon]: !!icon,
  });
    return (
      ReactDom.createPortal(
        <div className="Modal">
          <div
            className="Modal__cover"
            onClick={(evt) => {if(!noCancel) {handleClose(evt)}}}
            role="presentation"
          />

          <div className={modalContentCSS} style={{ overflow: overflowStyle }}>
            { handleClose && (
            <button
              type="button"
              className="Btn Btn--flat Modal__close padding--small "
              onClick={(evt) => handleClose(evt)}
            />
            )}
            <div className="Modal__container">
              { subheader && (
              <p className="Modal__pre-header">
                {subheader}
              </p>
              )}
              { header && (
              <>
                <h1 className="Modal__header">
                  <div className={`Icon Icon--${icon}`} />
                  {header}
                </h1>
                <hr />
              </>
              )}
              <div className="Modal__sub-container">
                {children}
              </div>
            </div>
          </div>
        </div>,
        document.getElementById('modal'),
      )
    );
}


export default Modal;
