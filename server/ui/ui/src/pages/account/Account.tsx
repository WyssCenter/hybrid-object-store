// vendor
import React, { FC, useContext, useState, useEffect } from 'react';
import ReactTooltip from 'react-tooltip';
// context
import AppContext from 'Src/AppContext';
// enivornment
import { get } from 'Environment/createEnvironment';
// components
import { SectionCard } from 'Components/card/index';
import AccountLabel from './label/AccountLabel';
import { PrimaryButton } from 'Components/button/index';
import PasswordModal from './modal/PasswordModal';
// css
import './Account.scss';


const Account: FC = () => {
  // context
  const { profile } = useContext(AppContext).user;
  const { groups } = profile;
  const groupsSeparated = groups.split(',').join(', ');
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [changePasswordIsSupported, setChangePasswordIsSupported] = useState(false);

  useEffect(() => {
    get(`password`, true)
    .then(response => response.json())
    .then(data => setChangePasswordIsSupported(data.changePasswordIsSupported));
  }, [])

  return (
    <div className="Account">
      <PasswordModal
        hideModal={() => setPasswordModalVisible(false)}
        isVisible={passwordModalVisible}
      />
      <div className="Account__body grid">
        <SectionCard>
          <div className="Account__profile">
            <AccountLabel
              label="Username"
              value={profile.nickname}
            />
            <AccountLabel
              label="Name"
              value={profile.name}
            />
            <AccountLabel
              label="Email"
              value={profile.email}
            />
            <AccountLabel
              label="Role"
              value={profile.role}
            />

            <AccountLabel
              label="Groups"
              value={groupsSeparated}
            />
            <h6 className="AccountLabel__h6">
              Password:
            </h6>
            <div
              className="Account__password"
            >
              <div
                data-tip="You cannot change your password, please contact your server admin."
                role="presentation"
                data-tip-disable={changePasswordIsSupported}
              >
                <PrimaryButton
                  click={() => setPasswordModalVisible(true)}
                  text="Change Password"
                  disabled={!changePasswordIsSupported}
                />
                <ReactTooltip
                  place="bottom"
                />
              </div>
            </div>
          </div>
        </SectionCard>
      </div>
    </div>
  )
}

export default Account;
