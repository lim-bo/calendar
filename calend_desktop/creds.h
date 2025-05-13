#ifndef CREDS_H
#define CREDS_H
#include <QString>

struct credentials {
    QString mail;
    QString pass;
};

struct credentials_reg : public credentials {
    QString f_name;
    QString s_name;
    QString t_name;
    QString department;
    QString pos;
};

#endif // CREDS_H
