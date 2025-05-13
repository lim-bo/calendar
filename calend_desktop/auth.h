#ifndef AUTH_H
#define AUTH_H
#include <QDialog>
#include <QVector>
#include <QLineEdit>
#include <QObject>
#include "client.h"
#include <QCloseEvent>
#include "cfg.h"
namespace Ui {
class auth;
}

class auth : public QDialog
{
    Q_OBJECT

public:
    explicit auth(bool* fl, QString* uid,  QWidget *parent = nullptr);
    ~auth();

private slots:
    void on_login_button_clicked();

    void on_checkBox_toggled(bool checked);

    void on_reg_button_clicked();

private:
    Ui::auth *ui;
    Client cli;
    bool* fl;
    QString* uid;
};

#endif // AUTH_H
