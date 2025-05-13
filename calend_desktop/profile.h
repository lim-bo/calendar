#ifndef PROFILE_H
#define PROFILE_H

#include <QDialog>
#include "cfg.h"
#include "client.h"
namespace Ui {
class profile;
}

class profile : public QDialog
{
    Q_OBJECT

public:
    explicit profile(QString uid, QWidget *parent = nullptr);
    ~profile();
    void loadUserData();

private slots:
    void on_pushButton_clicked();

private:
    Ui::profile *ui;
    Client cli;
    QString uid;
};

#endif // PROFILE_H
