import React, { useState } from "react";
import "./Header.css";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
} from "@mui/material";

function Header() {
  const [loginDialogOpen, setLoginDialogOpen] = useState(false);
  const [registrationDialogOpen, setRegisterDialogOpen] = useState(false);

  const handleLoginDialog = () => {
    setLoginDialogOpen(!loginDialogOpen);
  };

  const handleRegistrationDialog = () => {
    setRegisterDialogOpen(!registrationDialogOpen);
  };
  return (
    <>
      <div className="header">
        <div className="logo">NewBetting</div>
        <Button onClick={handleLoginDialog}>Prijava</Button>
        <Button onClick={handleRegistrationDialog}>Registracija</Button>
      </div>

      {/* login dialog */}
      <div>
        <Dialog open={loginDialogOpen} onClose={handleLoginDialog}>
          <DialogTitle>Prijava</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              id="name"
              label="Korisničko ime / email adresa"
              type="email"
              fullWidth
              variant="standard"
            />
            <TextField
              autoFocus
              margin="dense"
              id="password"
              label="Lozinka"
              type="password"
              fullWidth
              variant="standard"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleLoginDialog}>Odustani</Button>
            <Button onClick={null}>Prijavi se</Button>
          </DialogActions>
        </Dialog>
      </div>

      {/* registration dialog */}

      <div>
        <Dialog
          open={registrationDialogOpen}
          onClose={handleRegistrationDialog}
        >
          <DialogTitle>Registririraj se</DialogTitle>
          <DialogContent></DialogContent>
          <DialogActions>
            <Button onClick={handleRegistrationDialog}>Odustani</Button>
            <Button onClick={null}>Registriraj se</Button>
          </DialogActions>
        </Dialog>
      </div>
    </>
  );
}

export default Header;
