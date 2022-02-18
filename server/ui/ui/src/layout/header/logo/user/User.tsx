// vendor
import React,
{
  FC,
  useContext,
  useEffect,
  useRef,
  useState,
} from 'react';
import classNames from 'classnames';
import { Link } from 'react-router-dom';
// hooks
import { useEventListener } from 'Hooks/events/index';
// context
import AppContext from '../../../../AppContext';
// css
import { ReactComponent as UpArrow } from 'Images/icons/arrows/up-primary.svg';
import { ReactComponent as DownArrow } from 'Images/icons/arrows/down-primary.svg';
import './User.scss';

interface UserProps {
  [key: string]: any;
}


interface Props {
  send: any;
}

const User: FC<Props> = ({ send }: Props) => {
  // context
  const appContext = useContext(AppContext);
  // refs
  const dropdownButtonRef = useRef(null);
  // state
  const username = appContext.user && appContext.user.profile
    ? appContext.user.profile.nickname
    : 'loading';
  const [menuVisibility, setMenuVisibility] = useState(false);
  /**
  * Method triggers logout flow
  * @param {}
  * @fires {send}
  * @return {void}
  */
  const handleLogout = () => {
    send("LOGOUT");
  }

  /**
  * Method provides a way for child componts to update state
  * @param {Object} evt
  * @fires setMenuVisibility
  */
  const windowClickHandler = (event: Event) => {
    if (dropdownButtonRef.current && !dropdownButtonRef.current.contains(event.target)) {
      setMenuVisibility(false);
    } else {
      setMenuVisibility(!menuVisibility);
    }
  }

  // Add event listener using our hook
  useEventListener('click', windowClickHandler);


  // declare css here
  const userArrowCSS = classNames({
    'User__arrow': true,
    'User__arrow--up': menuVisibility,
    'User__arrow--down': !menuVisibility,
  });

  return (

    <div className="User">
      <button
        className="User__button User__button--user"
        onClick={() => setMenuVisibility(!menuVisibility)}
        ref={dropdownButtonRef}
        type="button"
      >
        <span className="User__h6">{username}</span>
        <div className={userArrowCSS}>
          {
            menuVisibility && (
              <UpArrow />
            )
          }
          {
            !menuVisibility && (
              <DownArrow />
            )
          }
        </div>
      </button>
      { menuVisibility &&
        <menu className="User__menu">
          <Link to="/account">
            <li className="User__li">

                <button
                  className="User__button User__button--flat"
                >
                Account
                </button>

            </li>
          </Link>

          <Link to="/groups">
            <li className="User__li">

                <button
                  className="User__button User__button--flat"
                >
                Groups
                </button>

            </li>
          </Link>


          <Link to="/tokens">
            <li className="User__li">

              <button
                className="User__button User__button--flat"
              >
              Tokens
              </button>

            </li>
          </Link>

          <li
            className="User__li"
            onClick={handleLogout}
          >
            <button
              className="User__button User__button--flat"
            >
              Logout
            </button>
          </li>
        </menu>
      }
    </div>
  )
}


export default User;
