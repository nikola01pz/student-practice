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
        <Button onClick={handleLoginDialog}>Login</Button>
        <Button onClick={handleRegistrationDialog}>Register</Button>
      </div>

      {/* login dialog */}
      <div>
        <Dialog open={loginDialogOpen} onClose={handleLoginDialog}>
          <DialogTitle>Login</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              id="name"
              label="E-mail"
              type="email"
              fullWidth
              variant="standard"
            />
            <TextField
              autoFocus
              margin="dense"
              id="password"
              label="Password"
              type="password"
              fullWidth
              variant="standard"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleLoginDialog}>Quit</Button>
            <Button onClick={null}>Login</Button>
          </DialogActions>
        </Dialog>
      </div>

      {/* registration dialog */}
      <div>
        <Dialog
          open={registrationDialogOpen}
          onClose={handleRegistrationDialog}
        >
          <DialogTitle>Register</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              id="email"
              label="E-mail"
              type="email"
              fullWidth
              variant="standard"
            />

            <TextField
              autoFocus
              margin="dense"
              id="username"
              label="Username"
              type="text"
              fullWidth
              variant="standard"
            />

            <TextField
              autoFocus
              margin="dense"
              id="password"
              label="Password"
              type="text"
              fullWidth
              variant="standard"
            />

            <TextField
              autoFocus
              margin="dense"
              id="name"
              label="Name"
              type="text"
              fullWidth
              variant="standard"
            />

            <TextField
              autoFocus
              margin="dense"
              id="surname"
              label="Surname"
              type="text"
              fullWidth
              variant="standard"
            />

            <TextField
              autoFocus
              margin="dense"
              id="dateOfBirth"
              type="date"
              fullWidth
              variant="standard"
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleRegistrationDialog}>Quit</Button>
            <Button onClick={null}>Register</Button>
          </DialogActions>
        </Dialog>
      </div>
    </>
  );
}

export default Header;
